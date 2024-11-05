package mongora

import (
	"go.mongodb.org/mongo-driver/bson"
	"net/url"
)

func AppendTextSearchTitleFilter(queryParams url.Values, filter bson.D) (bson.D, string) {
	val := queryParams.Get("title")
	if val != "" {
		return append(filter, bson.E{Key: "$text", Value: bson.D{{Key: "$search", Value: val}}}), val
	}

	return filter, val
}

func AppendTextSearchNameFilter(queryParams url.Values, filter bson.D) (bson.D, string) {
	val := queryParams.Get("name")
	if val != "" {
		return append(filter, bson.E{Key: "$text", Value: bson.D{{Key: "$search", Value: val}}}), val
	}

	return filter, val
}

// AppendOrFilter : Custom or structure to search the array of given fields and value to search as regex expression
func AppendRegexOrFilter(filter bson.D, fields []string, val string) (bson.D, string) {
	if val != "" {
		// Create an array to store the $or conditions
		orQueries := bson.A{}

		// Iterate over each field and add a regex condition for it
		for _, field := range fields {
			orQueries = append(orQueries, bson.D{{field, bson.D{{"$regex", val}, {"$options", "i"}}}})
		}

		// Add the $or array to the filter
		findQueries := bson.D{
			{"$or", orQueries},
		}

		return append(filter, findQueries...), val
	}

	return filter, val
}

func AppendFindByTitleFilter(queryParams url.Values, filter bson.D) (bson.D, string) {
	val := queryParams.Get("title")
	if val != "" {
		findQueries := bson.D{
			{"$or", bson.A{
				bson.D{{"title.en", bson.D{{"$regex", val}, {"$options", "i"}}}},
				bson.D{{"title.mm", bson.D{{"$regex", val}, {"$options", "i"}}}},
			}},
		}
		return append(filter, findQueries...), val
	}

	return filter, val
}

func AppendFindByTitleAndKeyWordsFilter(queryParams url.Values, filter bson.D) (bson.D, string) {
	val := queryParams.Get("title")
	if val != "" {
		findQueries := bson.D{
			{"$or", bson.A{
				bson.D{{"title.en", bson.D{{"$regex", val}, {"$options", "i"}}}},
				bson.D{{"title.mm", bson.D{{"$regex", val}, {"$options", "i"}}}},
				bson.D{{"keywords.en", bson.D{{"$regex", val}, {"$options", "i"}}}},
				bson.D{{"keywords.mm", bson.D{{"$regex", val}, {"$options", "i"}}}},
			}},
		}
		return append(filter, findQueries...), val
	}

	return filter, val
}

func AppendFindByName(queryParams url.Values, filter bson.D) (bson.D, string) {
	val := queryParams.Get("name")
	if val != "" {
		findQueries := bson.D{
			{"$or", bson.A{
				bson.D{{"name.en", bson.D{{"$regex", val}, {"$options", "i"}}}},
				bson.D{{"name.mm", bson.D{{"$regex", val}, {"$options", "i"}}}},
			}},
		}
		return append(filter, findQueries...), val
	}

	return filter, val
}

func AppendFindByNameAndKeyWordsFilter(queryParams url.Values, filter bson.D) (bson.D, string) {
	val := queryParams.Get("name")
	if val != "" {
		findQueries := bson.D{
			{"$or", bson.A{
				bson.D{{"name.en", bson.D{{"$regex", val}, {"$options", "i"}}}},
				bson.D{{"name.mm", bson.D{{"$regex", val}, {"$options", "i"}}}},
				bson.D{{"keywords.en", bson.D{{"$regex", val}, {"$options", "i"}}}},
				bson.D{{"keywords.mm", bson.D{{"$regex", val}, {"$options", "i"}}}},
			}},
		}
		return append(filter, findQueries...), val
	}

	return filter, val
}
