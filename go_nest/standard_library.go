package goNest

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// HashPassword : Hash the password
func HashPassword(password string) (string, error) {
	// Generate salt and hash the password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ParseString : Convert any to string
func ParseString(val any) string {
	if val == nil {
		return ""
	}

	return fmt.Sprintf("%v", val)
}

func ParseInt(val any) int {
	return parseIntFromString(ParseString(val))
}

// Can also return multiple returns using (int, error)
func parseIntFromString(val string) int {
	//log.Print("Parsing integer: " + val)
	newVal := 0
	if val != "" {
		// Declare variables for the result and error
		intVal, err := strconv.Atoi(val)
		if err != nil {
			//log.Printf("ParseInt Error: %s", err)
			newVal = 0
		} else {
			newVal = intVal
		}
	}

	// Return the result and the error~
	return newVal
}

func ParseBool(val string) bool {
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return false // Return default false on error
	}
	return boolVal
}

func ToString(val interface{}) string {
	return fmt.Sprintf("%d", val)
}

func HideStringPartially(strToHide string) string {
	// Calculate the start index of the middle three characters
	startIndex := (len(strToHide) - 3) / 2

	// Replace the middle three characters with '*'
	return strToHide[:startIndex] + "***" + strToHide[startIndex+3:]
}

func GetMongoTime() primitive.DateTime {
	return primitive.NewDateTimeFromTime(time.Now())
}

func MapToString(data map[string]interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func ArrayToString(array []interface{}) string {
	if len(array) == 0 {
		return "[]"
	}
	jsonData, _ := json.Marshal(array)
	return string(jsonData)
}

func StringToArray(val string) []interface{} {
	// Variable to hold the decoded JSON data
	var results []interface{}

	// Unmarshal JSON string to []interface{}
	err := json.Unmarshal([]byte(val), &results)
	if err != nil {
		return []interface{}{}
	}

	return results
}

func ArrayToPlainString(val []string, separator string) string {
	commaSeparated := strings.Join(val, separator)
	return commaSeparated
}

func BsonArrayToString(docs []bson.M) (string, error) {
	// Convert the []bson.M to JSON format
	if len(docs) <= 0 {
		return "[]", nil
	}

	jsonData, err := json.Marshal(docs)
	if err != nil {
		return "", fmt.Errorf("error converting to JSON: %v", err)
	}

	// Convert JSON data to a string and return
	return string(jsonData), nil
}

func BytesToString(jsonData []byte, isDocument bool) string {
	strResults := ""
	if utf8.Valid(jsonData) {
		//fmt.Sprintf("%v", jsonData)
		strResults = string(jsonData)
	}

	if strResults == "" || strResults == "null" {
		if isDocument {
			strResults = "{}"
		} else {
			strResults = "[]"
		}
	}

	return strResults
}

// ConvertJsonBytesToMap : Convert JSON byte data into readable format
func ConvertJsonBytesToMap(jsonData []byte) (map[string]interface{}, string) {
	// Unmarshal the JSON back into a map
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Sprintf("Error unmarshalling JSON Map: %v", err)
	}

	return data, ""
}

func ConvertJsonBytesToArray(jsonData []byte) ([]interface{}, string) {
	// Unmarshal the JSON into an array
	var data []interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Sprintf("Error unmarshalling JSON Array: %v", err)
	}

	return data, ""
}

func ConvertMapToBsonD(data map[string]interface{}) bson.D {
	var doc bson.D

	// Loop over each key-value pair in the map
	for key, value := range data {
		doc = append(doc, bson.E{Key: key, Value: value})
	}

	return doc
}

// ConvertBsonToByte : Convert BSON document to []byte
func ConvertBsonToByte(bsonData interface{}) ([]byte, error) {
	byteData, err := bson.Marshal(bsonData)
	if err != nil {
		return nil, err
	}

	return byteData, nil
}

func ConvertBsonToString(bsonData interface{}) (string, error) {
	//jsonBytes, err := bson.MarshalExtJSON(bsonData, true, true) // This one print extended bson document such as integer as $int
	jsonBytes, err := json.Marshal(bsonData)
	if err != nil {
		return "", err
	}

	// Convert JSON bytes to a string
	jsonString := string(jsonBytes)
	return jsonString, nil
}

func ConvertArrayToBsonD(data []interface{}) ([]bson.D, error) {
	var docs []bson.D

	for _, item := range data {
		switch v := item.(type) {
		case map[string]interface{}:
			// Convert map to bson.D
			var doc bson.D
			for key, value := range v {
				doc = append(doc, bson.E{Key: key, Value: value})
			}
			docs = append(docs, doc)

		case string:
			// Add string as a single-entry bson.D with a default key name
			doc := bson.D{{Key: "value", Value: v}}
			docs = append(docs, doc)

		default:
			return nil, fmt.Errorf("unsupported type in array: %T", item)
		}
	}

	return docs, nil
}

func ArrayToInterfaceArray(songIds []interface{}) ([]interface{}, error) {
	var converted []interface{}
	for _, id := range songIds {
		converted = append(converted, id) // Append the ID directly without wrapping in bson.D
	}
	return converted, nil
}

// ArrayToObjectIDs converts a slice of interface{} IDs to a slice of ObjectIDs.
func ArrayToObjectIDs(songIds []interface{}) ([]interface{}, error) {
	var converted []interface{}
	for _, id := range songIds {
		// Ensure the id is a string
		strID, ok := id.(string)
		if !ok {
			return nil, fmt.Errorf("item is not a string: %v", id)
		}

		// Convert the string ID to ObjectID
		objID, err := primitive.ObjectIDFromHex(strID)
		if err != nil {
			return nil, fmt.Errorf("invalid ObjectID format for ID %v: %w", id, err)
		}

		converted = append(converted, objID)
	}
	return converted, nil
}

func ConvertBsonDToMap(doc bson.D) map[string]interface{} {
	result := make(map[string]interface{})

	// Iterate over each element in the bson.D slice
	for _, elem := range doc {
		result[elem.Key] = elem.Value
	}

	return result
}

// ConvertBsonMToMap : bson.M is already a Map and don't need particular conversion
func ConvertBsonMToMap(doc bson.M) map[string]interface{} {
	// Simply return the bson.M as it is
	return doc
}

func ConvertStringToBytes(val string) ([]byte, error) {
	byteSlice := []byte(val)
	return byteSlice, nil
}

func ConvertStringToMap(val string) (bson.M, error) {
	// Define a bson.M to hold the JSON data
	var result bson.M

	// Convert JSON string to bson.M
	err := json.Unmarshal([]byte(val), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetMin : retrieves specific fields from a map-like structure based on a comma-separated list of field names
func GetMin(data interface{}, fields string) map[string]interface{} {
	var dataMap map[string]interface{}
	result := make(map[string]interface{})

	// Determine the type of data and convert it to map[string]interface{} if necessary
	switch v := data.(type) {
	case map[string]interface{}:
		dataMap = v
	case bson.M:
		dataMap = map[string]interface{}(v)
	default:
		fmt.Println("The provided data is not of a supported type (map[string]interface{} or bson.M)")
		return result
	}

	// Split the comma-separated string into individual field names
	fieldList := strings.Split(fields, ",")

	// Loop through each field and add it to the result map if it exists in dataMap
	for _, field := range fieldList {
		field = strings.TrimSpace(field) // Remove any leading/trailing whitespace
		if val, exists := dataMap[field]; exists {
			result[field] = val
		}
	}

	return result
}

// GetDateTimeFromStringField retrieves a date string from a map-like structure and parses it as time.Time
func GetDateTimeFromStringField(data interface{}, field string) time.Time {
	// Default time set to 50 years before now
	defaultTime := time.Now().Add(-50 * 8760 * time.Hour)
	var dataMap map[string]interface{}

	// Determine the type of data and convert it to map[string]interface{} if necessary
	switch v := data.(type) {
	case map[string]interface{}:
		dataMap = v
	case bson.M:
		dataMap = map[string]interface{}(v)
	default:
		fmt.Println("The provided data is not of a supported type (map[string]interface{} or bson.M)")
		return defaultTime
	}

	// Retrieve the date string from the map
	dateStr, ok := dataMap[field].(string)
	if !ok {
		fmt.Printf("Field %s not found or is not a string\n", field)
		return defaultTime
	}

	// Parse the date string in ISO 8601 format
	parsedDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		fmt.Printf("Error parsing date: %v\n", err)
		return defaultTime
	}

	return parsedDate
}

func GetStringFromField(data interface{}, key string) string {
	var dataMap map[string]interface{}

	switch v := data.(type) {
	case map[string]interface{}:
		dataMap = v
	case bson.M:
		dataMap = map[string]interface{}(v)
	default:
		fmt.Println("The provided data is not of a supported type (map[string]interface{}, bson.M, or primitive.M)")
		return ""
	}

	// Check if the key exists and is a string
	if val, exists := dataMap[key]; exists {
		if strVal, ok := val.(string); ok {
			return strVal
		} else {
			fmt.Printf("The key '%s' is present, but the value is not a string. Actual type: %T\n", key, val)
		}
	} else {
		fmt.Printf("The key '%s' does not exist in the data map.\n", key)
	}

	return ""
}

// GetBoolFromField retrieves a boolean value from a map-like structure based on the given key
func GetBoolFromField(data interface{}, key string) bool {
	var dataMap map[string]interface{}

	switch v := data.(type) {
	case map[string]interface{}:
		dataMap = v
	case bson.M:
		dataMap = map[string]interface{}(v)
	default:
		fmt.Println("The provided data is not of a supported type (map[string]interface{} or bson.M)")
		return false
	}

	// Check if the key exists and is a bool
	if val, exists := dataMap[key]; exists {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		} else {
			fmt.Printf("The key '%s' is present, but the value is not a bool. Actual type: %T\n", key, val)
		}
	}

	return false
}

// GetIntFromField retrieves an integer value from a map-like structure based on the given key
func GetIntFromField(data interface{}, key string) int {
	var dataMap map[string]interface{}

	switch v := data.(type) {
	case map[string]interface{}:
		dataMap = v
	case bson.M:
		dataMap = map[string]interface{}(v)
	default:
		fmt.Println("The provided data is not of a supported type (map[string]interface{} or bson.M)")
		return 0
	}

	// Check if the key exists and is of an integer type
	if val, exists := dataMap[key]; exists {
		switch v := val.(type) {
		case uint8, uint16, uint32, uint64:
			return int(v.(uint64)) // Convert unsigned integer types to int safely
		case int, int8, int16, int32, int64:
			return int(v.(int64)) // Convert signed integer types to int safely
		default:
			fmt.Printf("The key '%s' is present, but the value is not an integer. Actual type: %T\n", key, val)
		}
	}

	return 0
}

// GetFloatFromField retrieves a float64 value from a map-like structure based on the given key
func GetFloatFromField(data interface{}, key string) float64 {
	var dataMap map[string]interface{}

	switch v := data.(type) {
	case map[string]interface{}:
		dataMap = v
	case bson.M:
		dataMap = map[string]interface{}(v)
	default:
		fmt.Println("The provided data is not of a supported type (map[string]interface{} or bson.M)")
		return 0.0
	}

	// Check if the key exists and is a float
	if val, exists := dataMap[key]; exists {
		switch v := val.(type) {
		case float32:
			return float64(v) // Convert float32 to float64
		case float64:
			return v
		default:
			fmt.Printf("The key '%s' is present, but the value is not a float. Actual type: %T\n", key, val)
		}
	}

	return 0.0
}

// GetObjectFromField retrieves a nested map from a map-like structure based on the given key
func GetObjectFromField(data interface{}, key string) map[string]interface{} {
	var dataMap map[string]interface{}

	switch v := data.(type) {
	case map[string]interface{}:
		dataMap = v
	case bson.M:
		dataMap = map[string]interface{}(v)
	default:
		fmt.Println("The provided data is not of a supported type (map[string]interface{} or bson.M)")
		return map[string]interface{}{}
	}

	// Check if the key exists and is a map-like structure
	if val, exists := dataMap[key]; exists {
		switch v := val.(type) {
		case map[string]interface{}:
			return v
		case bson.M:
			return map[string]interface{}(v) // Convert bson.M to map[string]interface{}
		default:
			fmt.Printf("The key '%s' is present, but the value is not a map. Actual type: %T\n", key, val)
		}
	}

	return map[string]interface{}{}
}

// GetArrayFromField retrieves an array (slice) from a map-like structure based on the given key
func GetArrayFromField(data interface{}, key string) []interface{} {
	var dataMap map[string]interface{}

	switch v := data.(type) {
	case map[string]interface{}:
		dataMap = v
	case bson.M:
		dataMap = map[string]interface{}(v)
	default:
		fmt.Println("The provided data is not of a supported type (map[string]interface{} or bson.M)")
		// Return an empty slice if data is not of the expected type
		return []interface{}{}
	}

	// Check if the key exists and is of an array type
	if val, exists := dataMap[key]; exists {
		switch v := val.(type) {
		case []interface{}:
			return v
		case primitive.A:
			return []interface{}(v) // Convert primitive.A to []interface{}
		default:
			fmt.Printf("The key '%s' is present, but the value is not an array. Actual type: %T\n", key, val)
		}
	} else {
		fmt.Printf("The key '%s' does not exist in the data map.\n", key)
	}

	// Return an empty slice if key not found or value isn't an array type
	return []interface{}{}
}

// GetDateTimeFromField retrieves a primitive.DateTime from a map-like structure based on the given key
func GetDateTimeFromField(data interface{}, key string) primitive.DateTime {
	var dataMap map[string]interface{}

	switch v := data.(type) {
	case map[string]interface{}:
		dataMap = v
	case bson.M:
		dataMap = map[string]interface{}(v)
	default:
		fmt.Println("The provided data is not of a supported type (map[string]interface{} or bson.M)")
		// Return a default DateTime value if data is not of the expected type
		return primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Second))
	}

	// Check if the key exists and is of type primitive.DateTime
	if val, exists := dataMap[key]; exists {
		if dateTimeVal, ok := val.(primitive.DateTime); ok {
			return dateTimeVal
		} else {
			fmt.Printf("The key '%s' is present, but the value is not a primitive.DateTime. Actual type: %T\n", key, val)
		}
	} else {
		fmt.Printf("The key '%s' does not exist in the data map.\n", key)
	}

	// Return a default DateTime value if key not found or not a DateTime type
	return primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Second))
}

//
//if val, ok := data[key].(string); ok {
//	return val
//} else {
//	//fmt.Println(key + " field not found or not a string")
//}

//func GetBoolFromField(data map[string]interface{}, key string) bool {
//	if val, ok := data[key].(bool); ok {
//		return val
//	} else {
//		//fmt.Println(key + " field not found or not a boolean")
//	}
//
//	return false
//}

//func GetIntFromField(data map[string]interface{}, key string) int {
//	// Check if the value is an int
//	if value, ok := data[key]; ok {
//		switch v := value.(type) {
//		case uint:
//		case uint8:
//		case uint16:
//		case uint32:
//		case uint64:
//		case int:
//		case int8:
//		case int16:
//		case int32:
//		case int64:
//			return int(v)
//		default:
//			fmt.Println("Value is not an integer type")
//		}
//	} else {
//		fmt.Println("Key not found in map")
//	}
//
//	// Key not found or unsupported type
//	// fmt.Println(key + " field not found or not an integer/float")
//	return 0
//}

//func GetFloatFromField(data map[string]interface{}, key string) float64 {
//	if value, ok := data[key]; ok {
//		switch v := value.(type) {
//		case float32:
//		case float64:
//			return v
//		default:
//			fmt.Println("Value is not an integer type")
//		}
//	} else {
//		fmt.Println("Key not found in map")
//	}
//
//	return 0.0
//}
//
//func GetObjectFromField(data map[string]interface{}, key string) map[string]interface{} {
//	if value, ok := data[key]; ok {
//		switch v := value.(type) {
//		case map[string]interface{}:
//			return v
//		case primitive.M:
//			return v
//		default:
//			fmt.Println("Value is not an map type")
//		}
//	} else {
//		fmt.Println("Key not found in map")
//	}
//
//	// Return an empty map if key is not found or not the expected type
//	return map[string]interface{}{}
//}

//func GetArrayFromField(data map[string]interface{}, key string) []interface{} {
//	if value, ok := data[key]; ok {
//		switch v := value.(type) {
//		case []interface{}:
//			return v
//		case primitive.A:
//			return v
//		default:
//			fmt.Println("Value is not an map type")
//		}
//	} else {
//		fmt.Println("Key not found in map")
//	}
//
//	return []interface{}{}
//}
//
//func GetDateTimeFromField(data map[string]interface{}, key string) primitive.DateTime {
//	if value, ok := data[key]; ok {
//		switch v := value.(type) {
//		case primitive.DateTime:
//			return v
//		default:
//			fmt.Println("Value is not an map type")
//		}
//	} else {
//		fmt.Println("Key not found in map")
//	}
//
//	return primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Second))
//}

func ConvertDateTimeToString(dateTime primitive.DateTime) string {
	isoString := dateTime.Time().Format(time.RFC3339)
	return isoString
}

// GenerateRandomUUID : Generate a unique request ID (similar to crypto.randomUUID)
func GenerateRandomUUID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(bytes)
}

func GetPublicIP() string {
	ipAddress, err := getMachineIP(true)
	if err != nil {
		return ""
	}

	return ipAddress
}

func GetLocalIP() string {
	ipAddress, err := getMachineIP(false)
	if err != nil {
		return ""
	}

	return ipAddress
}

func getMachineIP(isPublicIP bool) (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		// Check if the address is an IP address and is not a loopback
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ipAddress := ipNet.IP.String()
			if isPublicIP {
				return ipAddress, nil
			} else {
				// Ensure that the IP is in a private range
				if isPrivateIP(ipAddress) {
					return ipAddress, nil
				}
			}

		}
	}
	return "", fmt.Errorf("no private local IP address found")
}

func isPrivateIP(ip string) bool {
	return strings.HasPrefix(ip, "10.") ||
		strings.HasPrefix(ip, "172.16.") ||
		strings.HasPrefix(ip, "172.17.") ||
		strings.HasPrefix(ip, "172.18.") ||
		strings.HasPrefix(ip, "172.19.") ||
		strings.HasPrefix(ip, "172.20.") ||
		strings.HasPrefix(ip, "172.21.") ||
		strings.HasPrefix(ip, "172.22.") ||
		strings.HasPrefix(ip, "172.23.") ||
		strings.HasPrefix(ip, "172.24.") ||
		strings.HasPrefix(ip, "172.25.") ||
		strings.HasPrefix(ip, "172.26.") ||
		strings.HasPrefix(ip, "172.27.") ||
		strings.HasPrefix(ip, "172.28.") ||
		strings.HasPrefix(ip, "172.29.") ||
		strings.HasPrefix(ip, "172.30.") ||
		strings.HasPrefix(ip, "172.31.") ||
		strings.HasPrefix(ip, "192.168.")
}

// generateChecksum calculates the SHA-256 hash of the provided file content and returns it as a hex string.
func GenerateChecksum(file *os.File) (string, error) {
	// Initialize SHA-256 hasher
	hasher := sha256.New()

	// Read the file in chunks and update the hash
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	// Compute the final hash and convert it to a hex string
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
