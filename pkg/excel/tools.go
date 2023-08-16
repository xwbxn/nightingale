package excel

import (
	"fmt"
	"math/rand"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/gin-gonic/gin"
)

var (
	defaultSheetName = "Sheet1" //默认Sheet名称
	defaultHeight    = 25.0     //默认行高度
)

type lzExcelExport struct {
	file      *excelize.File
	sheetName string //可定义默认sheet名称
}

func NewMyExcel(theme string) *lzExcelExport {
	defaultSheetName = theme
	return &lzExcelExport{file: createFile(), sheetName: defaultSheetName}
}

//导出基本的表格
func (l *lzExcelExport) ExportToPath(params []map[string]string, data []map[string]interface{}, path string) (string, error) {
	l.export(params, data)
	name := createFileName()
	filePath := path + "/" + name
	err := l.file.SaveAs(filePath)
	return filePath, err
}

func ReadExce[T any](xlsx *excelize.File, ctx *ctx.Context) ([]T, error) {
	//根据名字获取cells的内容，返回的是一个[][]string
	rows := xlsx.GetRows(xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))
	//声明一个数组
	var entitys []T
	fields := reflect.ValueOf(new(T)).Elem()
	mapLit := make(map[int]string)
	for i, row := range rows {
		//去掉第一行是excel表头部分
		if i == 0 { //取得第一行的所有数据---execel表头
			for index, colCell := range row {
				mapLit[index] = colCell
			}
		} else {
			entity := new(T)
			g := reflect.ValueOf(entity).Elem()
			for index, colCell := range row {
				title := mapLit[index]
				for i := 0; i < fields.NumField(); i++ {
					fieldInfo := fields.Type().Field(i)

					cnTag, heardOk := fieldInfo.Tag.Lookup("cn")
					sourceTag, sourceOk := fieldInfo.Tag.Lookup("source")

					if heardOk && (cnTag == title) {
						var results []int64
						var isDB = false

						if sourceOk {
							sources := strings.Split(sourceTag, ",")
							m := make(map[string]string)
							for _, pair := range sources {
								kv := strings.Split(pair, "=")
								m[kv[0]] = strings.Trim(kv[1], " ")
							}
							if m["type"] == "table" {
								isDB = true
								session := models.DB(ctx)
								session.Table(m["table"]).Where(m["field"]+" = ?", colCell).Pluck("id", &results)
							} else if m["type"] == "option" {
								isDB = true
								currentValue := m["value"][1 : len(m["value"])-1]
								rangs := strings.Split(currentValue, ";")
								for idx := 0; idx < len(rangs); idx++ {
									if rangs[idx] == colCell {
										results = append(results, int64(idx))
										break
									}
								}
							} else if m["type"] == "date" {
								isDB = true
								timeLayout := "2006-01-02"
								times, _ := time.Parse(timeLayout, colCell)
								results = append(results, int64(times.Unix()))
							}
						}
						switch fieldType := fieldInfo.Type.Kind(); fieldType {
						case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
							{
								if isDB {
									if len(results) > 0 {
										g.FieldByName(fieldInfo.Name).SetInt(results[0])
									}
								} else {
									s1, _ := strconv.Atoi(colCell)
									g.FieldByName(fieldInfo.Name).SetInt(int64(s1))
								}
							}
						case reflect.String:
							g.FieldByName(fieldInfo.Name).SetString(colCell)
						case reflect.Bool:
							g.FieldByName(fieldInfo.Name).SetBool(colCell == "true")
						default:
							fmt.Printf("field type %s not support yet", fieldType)
						}
					}
				}

			}
			entitys = append(entitys, *entity)
		}

	}
	return entitys, nil

}

//导出到浏览器--本系统调用方法体。此处使用的gin框架 其他框架可自行修改ctx
func (l *lzExcelExport) ExportDataInfo(data []interface{}, tagName string, ctx *ctx.Context, gtx *gin.Context) {
	dataKey := make([]map[string]string, 0)    //表头
	datas := make([]map[string]interface{}, 0) //数据

	if data != nil && len(data) > 0 {
		t := reflect.TypeOf(data[0])
		v := reflect.ValueOf(data[0])
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		// 判断是否是结构体
		if v.Kind() != reflect.Struct {
			fmt.Println("it is not struct")
			return
		}
		for i := 0; i < t.NumField(); i++ {
			_, ok := t.Field(i).Tag.Lookup(tagName)
			if ok == true {

			}
		}
		for i := 0; i < t.NumField(); i++ {

			fieldInfo := t.Field(i)
			_, ok := fieldInfo.Tag.Lookup(tagName)
			if ok == true {
				var is_num string = "0"
				switch fieldType := fieldInfo.Type.Kind(); fieldType {
				case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
					is_num = "1"
				}
				dataKey = append(dataKey, map[string]string{
					"key":    fieldInfo.Name,
					"title":  string(fieldInfo.Tag.Get(tagName)),
					"width":  "20",
					"is_num": is_num,
				})
			}
		}

		for _, entity := range data {
			out := make(map[string]interface{})
			t := reflect.TypeOf(entity)
			v := reflect.ValueOf(entity)
			// 取出指针的值
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			// 判断是否是结构体
			if v.Kind() != reflect.Struct {
				fmt.Println("it is not struct")
			}

			for i := 0; i < t.NumField(); i++ {
				_, cnOk := t.Field(i).Tag.Lookup(tagName)
				sourceTag, sourceOk := t.Field(i).Tag.Lookup("source")
				if cnOk == true {

					if sourceOk {
						m := make(map[string]string)
						sources := strings.Split(sourceTag, ",")
						for _, pair := range sources {
							kv := strings.Split(pair, "=")
							m[kv[0]] = strings.Trim(kv[1], " ")
						}

						if m["type"] == "table" {
							session := models.DB(ctx)
							var results []string
							session.Table(m["table"]).Where(" id = ?", v.Field(i).Interface()).Pluck(m["field"], &results)
							if len(results) > 0 {
								out[t.Field(i).Name] = results[0]
							}
						} else if m["type"] == "option" {
							currentValue := m["value"][1 : len(m["value"])-1]
							rangs := strings.Split(currentValue, ";")
							for idx := 0; idx < len(rangs); idx++ {
								currentValue := fmt.Sprintf("%d", v.Field(i).Interface())
								if fmt.Sprintf("%d", idx) == currentValue {
									out[t.Field(i).Name] = rangs[idx]
									break
								}
							}
						} else {
							out[t.Field(i).Name] = v.Field(i).Interface()
						}

					} else {
						out[t.Field(i).Name] = v.Field(i).Interface()
					}

				}

			}
			datas = append(datas, out)
		}
	}

	l.export(dataKey, datas)
	buffer, _ := l.file.WriteToBuffer()
	//设置文件类型
	gtx.Header("Content-Type", "application/vnd.ms-excel;charset=utf8")
	//设置文件名称
	gtx.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(createFileName()))
	_, _ = gtx.Writer.Write(buffer.Bytes())
}

//导出到浏览器。此处使用的gin框架 其他框架可自行修改ctx
func (l *lzExcelExport) ExportTempletToWeb(data []interface{}, tagName string, optionTagName string, ctx *ctx.Context, c *gin.Context) {

	dataKey := make([]map[string]string, 0) //表头

	if data != nil && len(data) > 0 {
		t := reflect.TypeOf(data[0])
		v := reflect.ValueOf(data[0])
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		// 判断是否是结构体
		if v.Kind() != reflect.Struct {
			fmt.Println("it is not struct")
			return
		}
		for i := 0; i < t.NumField(); i++ {
			_, ok := t.Field(i).Tag.Lookup(tagName)
			if ok == true {

			}
		}
		for i := 0; i < t.NumField(); i++ {

			fieldInfo := t.Field(i)
			cnTag, cnOk := fieldInfo.Tag.Lookup(tagName)
			if cnOk == true {
				var is_num string = "0"
				switch fieldType := fieldInfo.Type.Kind(); fieldType {
				case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
					is_num = "1"
				}
				dataKey = append(dataKey, map[string]string{
					"key":    fieldInfo.Name,
					"title":  string(cnTag),
					"width":  "20",
					"is_num": is_num,
				})
			}
		}
	}

	l.exportTemplet(dataKey, data, ctx)
	buffer, _ := l.file.WriteToBuffer()
	//设置文件类型
	c.Header("Content-Type", "application/vnd.ms-excel;charset=utf8")
	//设置文件名称
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(createFileName()))
	_, _ = c.Writer.Write(buffer.Bytes())
}

//导出到浏览器。此处使用的gin框架 其他框架可自行修改ctx
func (l *lzExcelExport) ExportToWeb(params []map[string]string, data []map[string]interface{}, c *gin.Context) {
	l.export(params, data)
	buffer, _ := l.file.WriteToBuffer()
	//设置文件类型
	c.Header("Content-Type", "application/vnd.ms-excel;charset=utf8")
	//设置文件名称
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(createFileName()))
	_, _ = c.Writer.Write(buffer.Bytes())
}

//设置首行
//调整头部列数超过26列，需要调整

func (l *lzExcelExport) writeTop(params []map[string]string) {
	topStyle, _ := l.file.NewStyle(`{"font":{"bold":true},"alignment":{"horizontal":"center","vertical":"center"}}`)
	var word = 'A'
	var num int = 1
	var prefixWord = ""
	//首行写入
	for _, conf := range params {
		title := conf["title"]
		width, _ := strconv.ParseFloat(conf["width"], 64)

		if num/26 == 2 {
			prefixWord = "AA"
		} else if num/26 == 1 {
			prefixWord = "A"
		}
		line := prefixWord + fmt.Sprintf("%c1", word)
		//设置标题
		l.file.SetCellValue(l.sheetName, line, title)
		//列宽
		l.file.SetColWidth(l.sheetName, prefixWord+fmt.Sprintf("%c", word), prefixWord+fmt.Sprintf("%c", word), width)
		//设置样式
		l.file.SetCellStyle(l.sheetName, line, line, topStyle)

		word++
		if num%26 == 0 {
			word = 'A'
		}
		num++
	}
}

//写入数据
func (l *lzExcelExport) writeData(params []map[string]string, data []map[string]interface{}) {
	lineStyle, _ := l.file.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center"}}`)
	lineStyleLeft, _ := l.file.NewStyle(`{"alignment":{"horizontal":"left","vertical":"center"}}`)
	//数据写入
	var j = 2 //数据开始行数
	for i, val := range data {
		//设置行高
		l.file.SetRowHeight(l.sheetName, i+1, defaultHeight)
		//逐列写入
		var word = 'A'
		var num int = 1
		var prefixWord = ""
		for _, conf := range params {
			valKey := conf["key"]

			isNum := conf["is_num"]
			if num/26 == 2 {
				prefixWord = "AA"
			} else if num/26 == 1 {
				prefixWord = "A"
			}
			line := prefixWord + fmt.Sprintf("%c%v", word, j)
			//设置值
			if isNum != "0" {
				valNum := fmt.Sprintf("%v", val[valKey])

				if strings.HasSuffix(valKey, "At") {
					int64Num, _ := strconv.ParseInt(valNum, 10, 64)
					dataTimeStr := time.Unix(int64Num, 0).Format("2006-01-02")
					l.file.SetCellValue(l.sheetName, line, dataTimeStr)
				} else {
					l.file.SetCellValue(l.sheetName, line, valNum)
				}
				l.file.SetCellStyle(l.sheetName, line, line, lineStyle)
			} else {
				l.file.SetCellValue(l.sheetName, line, val[valKey])
				//设置样式
				l.file.SetCellStyle(l.sheetName, line, line, lineStyleLeft)
			}

			word++
			if num%26 == 0 {
				word = 'A'
			}
			num++
		}
		j++
	}
	//设置行高 尾行
	l.file.SetRowHeight(l.sheetName, len(data)+1, defaultHeight)
}

//写入数据
func (l *lzExcelExport) writeTemplet(params []map[string]string, data []interface{}, ctx *ctx.Context) {
	lineStyle, _ := l.file.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center"}}`)
	lineStyleLeft, _ := l.file.NewStyle(`{"alignment":{"horizontal":"left","vertical":"center"}}`)
	//数据写入
	var j int = 2 //数据开始行数
	//设置行高
	l.file.SetRowHeight(l.sheetName, 1, defaultHeight)
	//逐列写入
	var word = 'A'
	var num int = 1
	var prefixWord = ""
	// for i, _ := range data {
	t := reflect.TypeOf(data[0])
	v := reflect.ValueOf(data[0])
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// 判断是否是结构体
	if v.Kind() != reflect.Struct {
		fmt.Println("it is not struct")
		return
	}
	for i := 0; i < t.NumField(); i++ {
		fieldInfo := t.Field(i)
		_, cnOK := fieldInfo.Tag.Lookup("cn")
		sourceTag, sourceOk := fieldInfo.Tag.Lookup("source")
		if num/26 == 2 {
			prefixWord = "AA"
		} else if num/26 == 1 {
			prefixWord = "A"
		}
		if cnOK == true {
			line := prefixWord + fmt.Sprintf("%c%v", word, j)
			var isNum string = "0"
			switch fieldType := fieldInfo.Type.Kind(); fieldType {
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
				isNum = "1"
			}
			if sourceOk {
				var values []string
				dv := excelize.NewDataValidation(true)
				dv.SetSqref(fmt.Sprintf("%c%v:%c1048576", word, j, word))
				m := make(map[string]string)
				sources := strings.Split(sourceTag, ",")
				for _, pair := range sources {
					kv := strings.Split(pair, "=")
					m[kv[0]] = strings.Trim(kv[1], " ")
				}

				if m["type"] == "table" {
					session := models.DB(ctx)
					var results []string
					session.Table(m["table"]).Pluck(m["field"], &results)
					values = append(values, results...)
				} else if m["type"] == "range" {
					currentValue := m["value"][1 : len(m["value"])-1]
					rangs := strings.Split(currentValue, ";")
					for _, pair := range rangs {
						starts := strings.Split(pair, "-")
						start, _ := strconv.Atoi(starts[0])
						var end int = start
						if len(starts) > 1 {
							val2, _ := strconv.Atoi(starts[1])
							end = val2
						}
						for m := start; m <= end; m++ {
							values = append(values, strconv.Itoa(m))
						}
					}

				} else if m["type"] == "option" {
					currentValue := m["value"][1 : len(m["value"])-1]
					rangs := strings.Split(currentValue, ";")
					for _, pair := range rangs {
						values = append(values, pair)
					}
				}
				dv.SetDropList(values)
				l.file.AddDataValidation(l.sheetName, dv)
			}
			if isNum != "0" {
				l.file.SetCellStyle(l.sheetName, line, line, lineStyle)
			} else {
				l.file.SetCellStyle(l.sheetName, line, line, lineStyleLeft)
			}
			word++
			if num%26 == 0 {
				word = 'A'
			}
			num++
		}
	}

	//设置行高 尾行
	l.file.SetRowHeight(l.sheetName, len(data)+1, defaultHeight)
}

func (l *lzExcelExport) export(params []map[string]string, data []map[string]interface{}) {
	l.writeTop(params)
	l.writeData(params, data)
}

func (l *lzExcelExport) exportTemplet(params []map[string]string, data []interface{}, ctx *ctx.Context) {
	l.writeTop(params)
	l.writeTemplet(params, data, ctx)
}

func createFile() *excelize.File {
	f := excelize.NewFile()
	// 创建一个默认工作表
	sheetName := defaultSheetName
	index := f.NewSheet(sheetName)
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	return f
}

func createFileName() string {
	name := time.Now().Format("2006-01-02-15-04-05")
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("excle-%v-%v.xlsx", name, rand.Int63n(time.Now().Unix()))
}

// Letter 遍历a-z
func Letter(length int) []string {
	var str []string
	for i := 0; i < length; i++ {
		str = append(str, string(rune('A'+i)))
	}
	return str
}

//excel导出(数据源为Struct) []interface{}
func (l *lzExcelExport) ExportExcelByStruct(theme string, titleList []string, data []interface{}, fileName string, sheetName string, c *gin.Context) error {
	l.file.SetSheetName(theme, sheetName)
	header := make([]string, 0)
	for _, v := range titleList {
		header = append(header, v)
	}
	rowStyleID, _ := l.file.NewStyle(`{"font":{"color":"#666666","size":13,"family":"arial"},"alignment":{"vertical":"center","horizontal":"center"}}`)

	l.file.SetSheetRow(sheetName, "A1", &header)

	l.file.SetRowHeight(theme, 1, 30)
	length := len(titleList)
	headStyle := Letter(length)
	var lastRow string
	var widthRow string
	for k, v := range headStyle {
		if k == length-1 {

			lastRow = fmt.Sprintf("%s1", v)
			widthRow = v
		}
	}
	l.file.SetColWidth(sheetName, "A", widthRow, 30)
	rowNum := 1
	for _, v := range data {

		t := reflect.TypeOf(v)
		fmt.Print("--ttt--", t.NumField())
		value := reflect.ValueOf(v)
		row := make([]interface {
		}, 0)
		for l := 0; l < t.NumField(); l++ {

			val := value.Field(l).Interface()
			row = append(row, val)
		}
		rowNum++
		l.file.SetCellStyle(sheetName, fmt.Sprintf("A%d", rowNum), fmt.Sprintf("%s", lastRow), rowStyleID)

	}
	disposition := fmt.Sprintf("attachment; filename=%s.xlsx", url.QueryEscape(fileName))
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", disposition)
	c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	return l.file.Write(c.Writer)
}
