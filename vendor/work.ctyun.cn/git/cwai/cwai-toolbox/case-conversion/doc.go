// This package offers the conversion between smallCamelCase, BigCamelCase
// and snake_case.
//
// The CamelToSnake is used to convert either smallCamelCase or BigCamelCase
// lossless to snake_case string. There is no exceptions cases.
//
// However, be careful when using SnakeToCamelSmall and SnakeToCamelBig to
// convert them back. Your special consequent-upper-case phrases may be lost
// during the conversion. For example:
//
//	columnNameJSON := "dataPipelineID"
//	columnNameDatabase := CamelToSnake(columnNameJSON)
//	println(columnNameDatabase)
//
// This will print "data_pipeline_id" without problem. However, if you try to
// use SnakeToCamelSmall to get it back, you may not get what you want:
//
//	columnNameDatabase := "data_pipeline_id"
//	columnNameJSON := SnakeToCamelSmall(columnNameDatabase)
//	println(columnNameJSON)
//
// This will print "dataPipelineId" instead of original "dataPipelineID". The
// special consequent-upper-case phrase "ID" is lost during conversion.
package case_conversion
