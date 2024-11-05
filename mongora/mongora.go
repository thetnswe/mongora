package mongora

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	goNest "github.com/thetnswe/mongora/go_nest"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// BodyValidate : Create a new global validator for request body
var BodyValidate *validator.Validate = validator.New()

// ////////////////
// // Queries ////
// ////////////////

func InsertOne(collection *mongo.Collection, reqBody interface{}) (primitive.ObjectID, error) {
	// Set a context with a timeout for the insert operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert the record
	recordId, err := collection.InsertOne(ctx, reqBody)
	if err != nil {
		return primitive.NilObjectID, err
	}

	//jsonData, err := json.Marshal(recordData)
	//if err != nil {
	//	return nil, fmt.Sprintf("Error converting BSON to JSON: %v", err)
	//}

	return recordId.InsertedID.(primitive.ObjectID), nil
}

func DeleteOne(collection *mongo.Collection, filter interface{}) (bool, error) {
	// Set a context with a timeout for the insert operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert the record
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}

	return true, nil
}

func FindOneAndUpdate(collection *mongo.Collection, filter interface{}, update interface{}) (bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result bson.M

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After) // Return the document after update
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func FindOneAndDelete(collection *mongo.Collection, filter interface{}) (bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result bson.M
	err := collection.FindOneAndDelete(ctx, filter).Decode(&result)

	if err != nil {
		return nil, err
		//if errors.Is(err, mongo.ErrNoDocuments) {
		//	return nil, "No document found"
		//} else {
		//	return nil, fmt.Sprintf("Fatal Error: %v", err)
		//}
	}

	result = bson.M{
		"message": "Record deleted",
	}

	return result, nil
}

func FindByIdOrSlug(collection *mongo.Collection, id string) (bson.M, error) {
	if IsValidObjectID(id) {
		return FindById(collection, id)
	}

	filter := bson.D{{"slug", id}}
	document, err := FindOne(collection, filter)
	return document, err
}

func FindById(collection *mongo.Collection, id string) (bson.M, error) {
	if !IsValidObjectID(id) {
		return nil, errors.New("Invalid Object ID")
	}

	//Find By ID
	objectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectID}

	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func FindOne(collection *mongo.Collection, filter interface{}) (bson.M, error) {
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ProcessFindQuery : Find Query processing
func Find(ctx context.Context, collection *mongo.Collection, filter interface{}) ([]bson.M, error) {
	results, err := FindWithAddonFields(ctx, collection, filter, "")
	return results, err
}

func FindWithAddonFields(ctx context.Context, collection *mongo.Collection, filter interface{}, addonFields string) ([]bson.M, error) {
	opts := buildOptionsForQuery(ctx, addonFields)
	cursor, err := collection.Find(ctx, filter, opts)

	if err != nil {
		return nil, err
	}

	//This function gets call after everything is completed, in this case, it gets called once the query is completed
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		//log.Println("Closing cursor")
		err := cursor.Close(ctx)
		if err != nil {
			_ = fmt.Sprintf("%v", err)
		}
	}(cursor, ctx)

	// Create a slice to hold the results
	var results []bson.M

	// Iterate over the cursor and decode each document to parse as JSON data later
	if err := cursor.All(ctx, &results); err != nil {
		return nil, errors.New(fmt.Sprintf("Error decoding documents: %v", err))
	}

	//log.Println("Query completed")
	return results, nil
}

// /////////////////
// // Utilities ////
// /////////////////
// GetIdFromModel : Return a string representation of ObjectId or String ID
func GetIdFromModel(model bson.M) (string, error) {
	idValue, exists := model["_id"]
	if !exists {
		return "", errors.New("missing _id field")
	}

	switch id := idValue.(type) {
	case primitive.ObjectID:
		// If _id is an ObjectID, convert to string
		strId := id.Hex() // Use Hex() method directly for ObjectID
		return strId, nil
	case string:
		// If _id is already a string, return it as-is
		return id, nil
	default:
		return "", errors.New("unsupported _id type")
	}
}

//func GetIdFromModel(model bson.M) (string, error) {
//	modelId, err := model["_id"].(primitive.ObjectID)
//	if !err {
//		//Also support string
//		modelId, err = model["_id"].(string)
//		if !err {
//			return "", errors.New("invalid ObjectID")
//		}
//	}
//
//	//Extract from quotes since .Hex() is not working for now
//	strId := modelId.String()
//	start := strings.Index(strId, "\"")
//	end := strings.LastIndex(strId, "\"")
//
//	if start == -1 || end == -1 || start == end {
//		log.Println("Invalid quotes", strId)
//		return "", errors.New("invalid Quotes")
//	}
//	strId = strId[start+1 : end]
//	return strId, nil
//
//}

func IsValidObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

func StringToObjectId(id string) (primitive.ObjectID, error) {
	strObjectId, err := primitive.ObjectIDFromHex(id)
	return strObjectId, err
}

// ObjectIdToString : Convert ObjectID to string
func ObjectIdToString(objectId primitive.ObjectID) string {
	strResult := objectId.Hex()
	return strResult
}

func GetSortString(ctx context.Context) string {
	orderBy := goNest.GetCtxStringValue(ctx, "order_by")
	order := goNest.GetCtxStringValue(ctx, "order")
	pageIndex := goNest.GetCtxStringValue(ctx, "page_index")
	countPerPage := goNest.GetCtxStringValue(ctx, "count_per_page")
	projection := goNest.GetCtxStringValue(ctx, "projection")
	projection = strings.ReplaceAll(projection, " ", "")

	return fmt.Sprintf("order_by=%s&order=%s&page_index=%s&count_per_page=%s&fields=%s", orderBy, order, pageIndex, countPerPage, projection)
}

func AppendFilter(queryParams url.Values, filter bson.D, key string) (bson.D, string) {
	val := queryParams.Get(key)
	if val != "" {
		return append(filter, bson.E{Key: key, Value: val}), val
	}

	return filter, val
}

func AppendBoolFilter(queryParams url.Values, filter bson.D, key string) (bson.D, string) {
	val := queryParams.Get(key)
	if val != "" {
		boolVal := goNest.ParseBool(val)
		return append(filter, bson.E{Key: key, Value: boolVal}), val
	}
	return filter, val
}

func ValidateRequestBody(req *http.Request, model any) (bool, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return false, errors.New("Unable to read request body")
	}

	// Unmarshal JSON if the body is expected to be JSON
	//var requestData karaoke.SongViewRequestData
	if err := json.Unmarshal(body, &model); err != nil {
		return false, errors.New("Request Body: Invalid JSON format")
	}

	// Validate the struct
	if err := BodyValidate.Struct(model); err != nil {
		return false, errors.New(fmt.Sprintf("Request Body validation error: %s", err))
	}

	return true, nil
}

// GenerateSearchTokens : Tokenized the given text so that it can search between spaces as well
func GenerateSearchTokens(text string) []string {
	// Split the text into words
	words := strings.Fields(text)
	var tokens []string

	// Generate all possible contiguous word sequences (n-grams)
	for start := 0; start < len(words); start++ {
		for end := start + 1; end <= len(words); end++ {
			token := strings.Join(words[start:end], " ")
			tokens = append(tokens, token)
		}
	}

	return tokens
}

// /////////////////////////////////
// // Private Utility Functions ////
// /////////////////////////////////
func buildOptionsForQuery(ctx context.Context, addonFields string) *options.FindOptions {
	opts := options.Find()

	// Define the projection
	fields := goNest.GetCtxStringValue(ctx, "fields")
	if addonFields != "" {
		if !(strings.Contains(fields, addonFields)) {
			fields = fields + "," + addonFields
		}
	}

	if fields != "" {
		projection, _ := parseProjection(fields)

		//TODO: Add multiple sort options for better searching and sorting
		opts = opts.SetProjection(projection)
	}

	//Check for other fields
	orderBy := goNest.GetCtxStringValue(ctx, "order_by")
	order := goNest.GetCtxStringValue(ctx, "order")

	sortOrder := bson.D{}
	if orderBy != "" && order != "" {
		orderValue := 1
		if order == "desc" {
			orderValue = -1
		}
		sortOrder = append(sortOrder, bson.E{Key: orderBy, Value: orderValue})

		opts.SetSort(sortOrder)
	}

	skip := int64(goNest.GetCtxIntValue(ctx, "skip"))
	limit := int64(goNest.GetCtxIntValue(ctx, "limit"))
	if skip >= 0 && limit > 0 {
		opts = opts.SetSkip(skip).SetLimit(limit)
	} else {
		skip = 0
		limit = 10
	}

	//log.Println(projection)
	//log.Println(sortOrder)
	//log.Println(skip)
	//log.Println(limit)

	return opts
}

func parseProjection(fields string) (bson.D, error) {
	if fields == "" {
		return bson.D{}, errors.New("fields string cannot be empty")
	}

	//Remove commas from the start and end character of the string
	fields = strings.Trim(fields, ",")

	parts := strings.Split(fields, ",")
	var projectionType int
	projection := bson.D{}

	for i, part := range parts {
		part = strings.TrimSpace(part)

		if part != "" {
			// Determine projection type based on prefix
			if strings.HasPrefix(part, "-") {
				if i == 0 {
					projectionType = 0
				} else if projectionType != 0 {
					return nil, errors.New("mix of inclusion and exclusion fields is not allowed")
				}
				part = strings.TrimPrefix(part, "-")
			} else {
				if i == 0 {
					projectionType = 1
				} else if projectionType != 1 {
					return nil, errors.New("mix of inclusion and exclusion fields is not allowed")
				}
			}

			// Append to projection
			projection = append(projection, bson.E{Key: part, Value: projectionType})
		}

	}

	return projection, nil
}

// /////////////////
// // Indexings ////
// /////////////////

// DropIndex : Drop index for the collection
func DropIndex(collection *mongo.Collection, indexName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := IndexExists(collection, indexName)
	if err != nil {
		log.Fatalf("Error checking index existence: %v", err)
	}

	//Only drop the index if the name exist
	if exists {
		// Drop the existing index
		_, err := collection.Indexes().DropOne(ctx, indexName)
		if err != nil {
			return err
		}
	} else {
		//No index exist
		//log.Printf("Index %s does not exist. Skipping drop operation.", indexName)
	}

	return nil
}

func IndexExists(collection *mongo.Collection, indexName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retrieve all indexes
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	// Loop through the indexes to find the specified index name
	for cursor.Next(ctx) {
		var index bson.M
		if err := cursor.Decode(&index); err != nil {
			return false, err
		}

		// Check if the index name matches
		if index["name"] == indexName {
			return true, nil
		}
	}

	return false, cursor.Err()
}

// CreateSingleIndexes : Helper function to create single indexes for the collection
func CreateSingleIndexes(collection *mongo.Collection, indexNames []string) {
	//log.Println("Creating single indexes")
	var wg sync.WaitGroup
	for _, indexName := range indexNames {
		wg.Add(1) // Increment WaitGroup counter
		go func(indexName string) {
			defer wg.Done() // Decrement counter when goroutine completes
			err := CreateIndexWithFields(collection, indexName, map[string]interface{}{indexName: 1})
			if err != nil {
				_ = fmt.Sprintf("%v\n", err)
			}
		}(indexName) // Pass indexName as argument to avoid capture by reference
	}
	wg.Wait() // Block until all goroutines have completed
	//log.Println("Single indexes created")
}

// CreateSingleHashIndexes : Helper function to create single hash indexes for the collection
func CreateSingleHashIndexes(collection *mongo.Collection, indexNames []string) {
	//log.Println("Creating single indexes")
	var wg sync.WaitGroup
	for _, indexName := range indexNames {
		wg.Add(1) // Increment WaitGroup counter
		go func(indexName string) {
			defer wg.Done() // Decrement counter when goroutine completes
			err := CreateIndexWithFields(collection, indexName, map[string]interface{}{indexName: "hashed"})
			if err != nil {
				_ = fmt.Sprintf("%v\n", err)
			}
		}(indexName) // Pass indexName as argument to avoid capture by reference
	}
	wg.Wait() // Block until all goroutines have completed
	//log.Println("Single indexes created")
}

// CreateIndexWithFields : Helper function to create an index with specified name and fields
func CreateIndexWithFields(collection *mongo.Collection, indexName string, indexes map[string]interface{}) error {
	// Drop the index if it already exists
	err := DropIndex(collection, indexName)
	if err != nil {
		return err
	}

	// Convert the map to bson.D format for the Keys
	keys := bson.D{}
	for field, value := range indexes {
		keys = append(keys, bson.E{Key: field, Value: value})
	}

	// Create the IndexModel with the specified keys and options
	indexModel := mongo.IndexModel{
		Keys: keys,
		Options: options.Index().SetName(indexName).SetDefaultLanguage("english").
			SetLanguageOverride("custom_language"),
	}

	// Set a timeout context for the index creation operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create the index on the specified collection
	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	return err
}
