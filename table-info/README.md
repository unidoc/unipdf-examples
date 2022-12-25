# Table Info examples

The `TableInfo()` method provides the table related information of the textMark that is being called from.
Generally it returns two kinds of information about the textMark in relation to a table.
* A pointer to the `Table` where the text mark is found in. If it is outside a table then this value is `nil`.
* The cell coordinate of the table where the texMark is found at. 
This is given as `[][]int` of row, column values.

## What can we do with this information?
Getting the table information from the text side gives us an access to the text related information of the tables. Here are some sample usages of this method.
1. Text distribution </br>
The distribution of the text inside tables vs out side tables. </br>
The code at [table_info.go](./table_info.go) demonstrates how to calculate the distribution of the page content.
2. Partition extracted text </br>
This means dividing the text into inside table and outside table sections. This can help us know where the table are found in relation to the text on the page. 
The example at  [partition_text.go](./partition_text.go) show cases how to divide an extracted text into sections of inside table and outside table contents. 