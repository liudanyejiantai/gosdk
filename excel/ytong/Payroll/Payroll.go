package Payroll //工资打款表

import (
	"container/list"
	"errors"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/liudanyejiantai/gosdk/string_func"
	"github.com/liudanyejiantai/gosdk/ulog"
)

const (
	PAY_ROLL_COUMNS = 15
)

var (
	log *ulog.Ulog
)

func init() {
	log = ulog.NewULog("", "")
}

type Payroll struct {
	SerialNo       string //序号
	EntryName      string //项目名称
	Account        string //银行帐号
	FullNameOfBank string //开户行全称
	IdNumber       string //身份证号

	Area                    string //区域
	Name                    string //员工姓名
	GrossPay                string //税前工资
	CorporateSocialSecurity string //企业社保
	PersonalSocialSecurity  string //个人社保

	EnterpriseProvidentFund string //企业公积金
	PersonalProvidentFund   string //个人公积金
	PersonalIncomeTax       string //个税
	ActualDistribution      string //实发
	ServiceCharge           string //服务费
}

//将array数组转换成struct
func (s *Payroll) ConvertFromArray(arr []string) error {
	if len(arr) != PAY_ROLL_COUMNS {
		log.WriteLog(ulog.ERROR, "数组长度%d,%d不等，无法转换成结构体Payroll", len(arr), PAY_ROLL_COUMNS)
		return errors.New("数组长度不对，无法转换成结构体Payroll")
	}

	s.SerialNo = arr[0]
	s.EntryName = arr[1]
	s.Account = arr[2]
	s.FullNameOfBank = arr[3]
	s.IdNumber = arr[4]

	s.Area = arr[5]
	s.Name = arr[6]
	s.GrossPay = arr[7]
	s.CorporateSocialSecurity = arr[8]
	s.PersonalSocialSecurity = arr[9]

	s.EnterpriseProvidentFund = arr[10]
	s.PersonalProvidentFund = arr[11]
	s.PersonalIncomeTax = arr[12]
	s.ActualDistribution = arr[13]
	s.ServiceCharge = arr[14]
	return nil
}

//struct转成array
func (s *Payroll) ConvertToArray() []string {
	var arr []string
	arr[0] = s.SerialNo
	arr[1] = s.EntryName
	arr[2] = s.Account
	arr[3] = s.FullNameOfBank
	arr[4] = s.IdNumber

	arr[5] = s.Area
	arr[6] = s.Name
	arr[7] = s.GrossPay
	arr[8] = s.CorporateSocialSecurity
	arr[9] = s.PersonalSocialSecurity

	arr[10] = s.EnterpriseProvidentFund
	arr[11] = s.PersonalProvidentFund
	arr[12] = s.PersonalIncomeTax
	arr[13] = s.ActualDistribution
	arr[14] = s.ServiceCharge
	return arr
}

//导入社保打款模板,生成对应的SQL语句
//for e := ls.Front(); e != nil; e = e.Next() {
//		s := e.Value.(*Payroll)
//		fmt.Println("s", s)
//	}
func ImportPayroll(strExcelFile string) (*list.List, error) {
	ls := list.New()

	xlsx, err := excelize.OpenFile(strExcelFile)
	if err != nil {
		log.WriteLog(ulog.ERROR, "打开EXCEL文件[%s]失败，原因:%s", strExcelFile, err.Error())
		return ls, err
	}
	//先用第一行取出全部的列
	rows := xlsx.GetRows("Sheet1")
	var colmuns []string
	for _, colCell := range rows[0] {
		colCell = string_func.Trim(colCell)
		if colCell != "" {
			colmuns = append(colmuns, colCell)
		}
	}

	for _, row := range rows {
		var arr []string
		i := 0
		for _, colCell := range row {
			i++
			arr = append(arr, colCell)
			if i == len(colmuns) {
				break
			}

		}
		var s Payroll
		s.ConvertFromArray(arr)
		ls.PushBack(&s)
		//fmt.Println("s", s)
	}

	return ls, nil
}
