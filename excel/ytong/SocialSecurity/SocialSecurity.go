package SocialSecurity

import (
	"container/list"
	"errors"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/liudanyejiantai/gosdk/string_func"
	"github.com/liudanyejiantai/gosdk/ulog"
)

var (
	log *ulog.Ulog
)

func init() {
	log = ulog.NewULog("", "")
}

const (
	SOCIAL_SEC_COLUMNS = 29 //社保信息一共29列
)

//社保信息
type SocialSec struct {
	PaymentCompany                      string //代缴公司
	Area                                string //地区
	SeialNo                             string //序号
	Name                                string //姓名
	EntryName                           string //项目名称
	SocialSecTime                       string //社保起缴时间
	PayTime                             string //缴纳时间
	CompanySecuritySubtotal             string //企业社保小计
	PersonSecuritySubtotal              string //个人社保小计
	CompanyProvidentFundSubtotal        string //企业公积金小计
	PersonProvidentFundSubtotal         string //个人公积金小计
	Total                               string //合计
	ServiceCharge                       string //服务费
	ProvidentFundPaymentBase            string //公积金缴纳基数
	PensionPaymentBase                  string //养老金缴纳基数
	MedicalInsurancePaymentBase         string //医疗保险缴纳基数
	EnterpriseEndowmentInsuranceFund    string //企业养老保险基金
	EnterpriseMedicalInsurance          string //企业医疗保险
	UnemploymentInsuranceForEnterprises string //企业失业险
	IndustrialInjuryInsurance           string //企业工伤保险
	EnterpriseMaternityInsurance        string //企业生育保险
	SeriousIllnessInEnterprises         string //企业大病
	EnterpriseDisabilityGold            string //企业残障金
	AnnualUnion                         string //工会费
	PersonalEndowmentInsuranceFund      string //个人养老保险基金
	IndividualHealthInsurance           string //个人医疗保险
	PersonalUnemploymentInsurance       string //个人失业险
	PersonalDisease                     string //个人大病
	CardMakingFee                       string //制卡费
}

//将array数组转换成struct
func (s *SocialSec) ConvertFromArray(arr []string) error {
	if len(arr) != SOCIAL_SEC_COLUMNS {
		log.WriteLog(ulog.ERROR, "数组长度%d,%d不等，无法转换成结构体SocialSec", len(arr), SOCIAL_SEC_COLUMNS)
		return errors.New("数组长度不对，无法转换成结构体SocialSec")
	}

	s.PaymentCompany = arr[0]
	s.Area = arr[1]
	s.SeialNo = arr[2]
	s.Name = arr[3]
	s.EntryName = arr[4]
	s.SocialSecTime = arr[5]
	s.PayTime = arr[6]
	s.CompanySecuritySubtotal = arr[7]
	s.PersonSecuritySubtotal = arr[8]
	s.CompanyProvidentFundSubtotal = arr[9]
	s.PersonProvidentFundSubtotal = arr[10]
	s.Total = arr[11]
	s.ServiceCharge = arr[12]
	s.ProvidentFundPaymentBase = arr[13]
	s.PensionPaymentBase = arr[14]
	s.MedicalInsurancePaymentBase = arr[15]
	s.EnterpriseEndowmentInsuranceFund = arr[16]
	s.EnterpriseMedicalInsurance = arr[17]
	s.UnemploymentInsuranceForEnterprises = arr[18]
	s.IndustrialInjuryInsurance = arr[19]
	s.EnterpriseMaternityInsurance = arr[20]
	s.SeriousIllnessInEnterprises = arr[21]
	s.EnterpriseDisabilityGold = arr[22]
	s.AnnualUnion = arr[23]
	s.PersonalEndowmentInsuranceFund = arr[24]
	s.IndividualHealthInsurance = arr[25]
	s.PersonalUnemploymentInsurance = arr[26]
	s.PersonalDisease = arr[27]
	s.CardMakingFee = arr[28]
	return nil
}

//struct转成array
func (s *SocialSec) ConvertToArray() []string {
	var arr []string
	arr[0] = s.PaymentCompany
	arr[1] = s.Area
	arr[2] = s.SeialNo
	arr[3] = s.Name
	arr[4] = s.EntryName
	arr[5] = s.SocialSecTime
	arr[6] = s.PayTime
	arr[7] = s.CompanySecuritySubtotal
	arr[8] = s.PersonSecuritySubtotal
	arr[9] = s.CompanyProvidentFundSubtotal
	arr[10] = s.PersonProvidentFundSubtotal
	arr[11] = s.Total
	arr[12] = s.ServiceCharge
	arr[13] = s.ProvidentFundPaymentBase
	arr[14] = s.PensionPaymentBase
	arr[15] = s.MedicalInsurancePaymentBase
	arr[16] = s.EnterpriseEndowmentInsuranceFund
	arr[17] = s.EnterpriseMedicalInsurance
	arr[18] = s.UnemploymentInsuranceForEnterprises
	arr[19] = s.IndustrialInjuryInsurance
	arr[20] = s.EnterpriseMaternityInsurance
	arr[21] = s.SeriousIllnessInEnterprises
	arr[22] = s.EnterpriseDisabilityGold
	arr[23] = s.AnnualUnion
	arr[24] = s.PersonalEndowmentInsuranceFund
	arr[25] = s.IndividualHealthInsurance
	arr[26] = s.PersonalUnemploymentInsurance
	arr[27] = s.PersonalDisease
	arr[28] = s.CardMakingFee
	return arr
}

//导入社保打款模板,生成对应的SQL语句
//for e := ls.Front(); e != nil; e = e.Next() {
//		s := e.Value.(*SocialSec)
//		fmt.Println("s", s)
//	}
func ImportSocialSecurity(strExcelFile string) (*list.List, error) {
	ls := list.New()
	xlsx, err := excelize.OpenFile(strExcelFile)
	if err != nil {
		log.WriteLog(ulog.ERROR, "打开EXCEL文件[%s]失败，原因:%s", strExcelFile, err.Error())
		return ls, err
	}
	//sheet页名称叫 社保
	rows := xlsx.GetRows("社保")

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
		var s SocialSec
		s.ConvertFromArray(arr)
		//fmt.Println("s", s)
		ls.PushBack(&s)
	}

	return ls, nil
}
