package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/devplayg/ipas-mcs/objs"
	"github.com/devplayg/ipas-mcs/libs"
	"regexp"
	"time"
	"github.com/devplayg/ipas-server"
)

var RegexFoundRows = regexp.MustCompile(`(?i)SELECT(\s+)SQL_CALC_FOUND_ROWS`)

// GetLogForCharting 함수와 같이 수정해야함
func GetIpaslog(filter *objs.IpasFilter, member *objs.Member) ([]objs.IpasLog, int64, error) {
	var where string
	var rows []objs.IpasLog

	// 조건 설정
	args := make([]interface{}, 0)

	// 검색시간 설정
	startDate, _ := time.ParseInLocation(ipasserver.DateDefault, filter.StartDate+":00", member.Location)
	endDate, _ := time.ParseInLocation(ipasserver.DateDefault, filter.EndDate+":59", member.Location)
	args = append(args, startDate.UTC().Format(ipasserver.DateDefault), endDate.UTC().Format(ipasserver.DateDefault))

	if member.Position < objs.Administrator {
		where += " and group_id in (select asset_id from mbr_asset where member_id = ?)"
		args = append(args, member.MemberId)
	}

	if len(filter.OrgId) > 0 {
		where += fmt.Sprintf(" and org_id in (%s)", libs.JoinInt(filter.OrgId, ","))
	}

	if len(filter.GroupId) > 0 {
		where += fmt.Sprintf(" and group_id in (%s)", libs.JoinInt(filter.GroupId, ","))
	}

	if len(filter.EventType) > 0 {
		where += fmt.Sprintf(" and event_type in (%s)", libs.JoinInt(filter.EventType, ","))
	}

	// 장비 태크 검색
	if len(filter.TagPattern) > 0 {
		where += " and (equip_id like ? or targets like ?)"
		cond := "%"+filter.TagPattern+"%"
		args = append(args, cond, cond)
	}
	if len(filter.EquipId) > 0 {
		where += " and equip_id = ?"
		cond := filter.EquipId
		args = append(args, cond)
	}

	// 페이징 모드(고속/일반)
	if filter.FastPaging == "off" {
		filter.FoundRows = "SQL_CALC_FOUND_ROWS"
	}

	// Set query
	query := `
		SELECT 	%s date, org_id, group_id, event_type, session_id, equip_id, targets, latitude, longitude, speed
				, snr, usim, distance, ip, recv_date
		from log_ipas_event
		where date >= ? and date <= ? %s
		order by %s %s
		limit ?, ?
	`
	query = fmt.Sprintf(query, filter.FoundRows, where, filter.Sort, filter.Order)
	args = append(args, filter.Offset, filter.Limit)

	o := orm.NewOrm()
	o.Begin()
	defer o.Commit()
	total, err := o.Raw(query, args).QueryRows(&rows)
	if filter.FastPaging == "off" {
		if RegexFoundRows.MatchString(query) {
			dbResult := objs.NewDbResult()
			o.Raw("select FOUND_ROWS() total").QueryRow(dbResult)
			total = dbResult.Total
		}
	}
	return rows, total, err
}


func GetLogForCharting(filter *objs.IpasFilter, member *objs.Member) ([]objs.StatsByEventType, int64, error) {
	var where string
	var rows []objs.StatsByEventType

	// 조건 설정
	args := make([]interface{}, 0)
	args = append(args, filter.StartDate+":00", filter.EndDate+":59")

	if member.Position < objs.Administrator {
		where += " and group_id in (select asset_id from mbr_asset where member_id = ?)"
		args = append(args, member.MemberId)
	}

	if len(filter.OrgId) > 0 {
		where += fmt.Sprintf(" and org_id in (%s)", libs.JoinInt(filter.OrgId, ","))
	}

	if len(filter.GroupId) > 0 {
		where += fmt.Sprintf(" and group_id in (%s)", libs.JoinInt(filter.GroupId, ","))
	}

	if len(filter.EventType) > 0 {
		where += fmt.Sprintf(" and event_type in (%s)", libs.JoinInt(filter.EventType, ","))
	}

	// 장비 태크 검색
	if len(filter.TagPattern) > 0 {
		where += " and (equip_id like ? or targets like ?)"
		cond := "%"+filter.TagPattern+"%"
		args = append(args, cond, cond)
	}
	if len(filter.EquipId) > 0 {
		where += " and equip_id = ?"
		cond := filter.EquipId
		args = append(args, cond)
	}

	// Set query
	query := `
		select concat(substr(date, 1, 13), ':00:00') date, event_type, count(*) count
		from log_ipas_event
		where date >= ? and date <= ? %s
		group by substr(date, 1, 13), event_type
	`
	query = fmt.Sprintf(query, where)

	o := orm.NewOrm()
	total, err := o.Raw(query, args).QueryRows(&rows)
	return rows, total, err
}

//func executeQuery(query string, args []interface{}, obj interface{})  {
//	o := orm.NewOrm()
//	o.Begin()
//	defer o.Commit()
//	spew.Dump(args)
//	refresh
//	objs := make([]obj.ref, 0)
//	o.Raw(query, args).QueryRows(&obj)
////
////	spew.Dump(obj)
//	//if RegexFoundRows.MatchString(query) {
//	//	dbResult := objs.NewDbResult()
//	//	o.Raw("select FOUND_ROWS() total").QueryRow(dbResult)
//	//	total = dbResult.Total
//	//}
//	//return obj, total, err
//}

//func GetIpasLog(filter *FilterFiletranslog) ([]objs.IpasLog, int64, error) {
//	var where string
//	args := make([]interface{}, 0, 20)
//
//	// Fast paging
//	if filter.FastPaging != "on" {
//		filter.FoundRows = "SQL_CALC_FOUND_ROWS"
//	}
//
//	// Time
//	args = append(args, filter.StartDate+":00")
//	args = append(args, filter.EndDate+":59")
//
//	// Network
//	if filter.Networks != nil {
//		where += fmt.Sprintf(" and network_id in (%s)", StringJoin(filter.Networks, ","))
//	} else {
//		if filter.Groups != nil {
//			where += fmt.Sprintf(" and group_id in (%s)", StringJoin(filter.Groups, ","))
//		} else {
//			if filter.Sensors != nil {
//				where += fmt.Sprintf(" and sensor_id in (%s)", StringJoin(filter.Sensors, ","))
//			} else {
//				// N/A
//			}
//		}
//
//	}
//	// Src IP
//	if len(filter.SrcIpCidr) > 0 {
//		if strings.Contains(filter.SrcIpCidr, "/") { // CIDR
//			ip_min, ip_max, err := libs.GetCidrMinMax(filter.SrcIpCidr)
//			CheckError(err)
//			if err == nil {
//				where += " and src_ip between ? and ?"
//				args = append(args, libs.IpToInt32(ip_min), libs.IpToInt32(ip_max))
//			}
//		} else { // Single IP
//			ip := net.ParseIP(filter.SrcIpCidr)
//			if ip != nil {
//				where += " and src_ip = ?"
//				args = append(args, libs.IpToInt32(ip))
//			}
//		}
//	}
//
//	// Src port
//	if len(filter.SrcPortStart) > 0 && len(filter.SrcPortEnd) > 0 {
//		where += " and src_port between ? and ?"
//		port_from, _ := strconv.Atoi(filter.SrcPortStart)
//		port_to, _ := strconv.Atoi(filter.SrcPortEnd)
//		args = append(args, port_from, port_to)
//	} else {
//		if len(filter.SrcPortStart) > 0 {
//			where += " and src_port >= ?"
//			port, _ := strconv.Atoi(filter.SrcPortStart)
//			args = append(args, port)
//		} else if len(filter.SrcPortEnd) > 0 {
//			where += " and src_port <= ?"
//			port, _ := strconv.Atoi(filter.SrcPortEnd)
//			args = append(args, port)
//		}
//	}
//
//	// Dst IP
//	if len(filter.DstIpCidr) > 0 {
//		if strings.Contains(filter.DstIpCidr, "/") { // CIDR
//			ip_min, ip_max, err := libs.GetCidrMinMax(filter.DstIpCidr)
//			CheckError(err)
//			if err == nil {
//				where += " and Dst_ip between ? and ?"
//				args = append(args, libs.IpToInt32(ip_min), libs.IpToInt32(ip_max))
//			}
//		} else { // Single IP
//			ip := net.ParseIP(filter.DstIpCidr)
//			if ip != nil {
//				where += " and Dst_ip = ?"
//				args = append(args, libs.IpToInt32(ip))
//			}
//		}
//	}
//
//	// Dst port
//	if len(filter.DstPortStart) > 0 && len(filter.DstPortEnd) > 0 {
//		where += " and Dst_port between ? and ?"
//		port_from, _ := strconv.Atoi(filter.DstPortStart)
//		port_to, _ := strconv.Atoi(filter.DstPortEnd)
//		args = append(args, port_from, port_to)
//	} else {
//		if len(filter.DstPortStart) > 0 {
//			where += " and Dst_port >= ?"
//			port, _ := strconv.Atoi(filter.DstPortStart)
//			args = append(args, port)
//		} else if len(filter.DstPortEnd) > 0 {
//			where += " and Dst_port <= ?"
//			port, _ := strconv.Atoi(filter.DstPortEnd)
//			args = append(args, port)
//		}
//	}
//
//	// Paging
//	args = append(args, filter.Offset)
//	args = append(args, filter.Limit)
//
//	// Set query
//	query := `
//        SELECT %s *
//        from log_event_filetrans
//        where rdate >= ? and rdate <= ? %s order by %s %s limit ?, ?
//    `
//	query = fmt.Sprintf(query, filter.FoundRows, where, filter.Sort, filter.Order)
//
//	o := orm.NewOrm()
//	rows := []RsFiletransLog{}
//	total, err := o.Raw(query, args).QueryRows(&rows)
//	dbResult := DbResult{}
//	isMatch := RegexFoundRows.MatchString(query)
//	if isMatch {
//		o.Raw("select FOUND_ROWS() total").QueryRow(&dbResult)
//		total = dbResult.Total
//
//	}
//
//	return &rows, total, err
//}
//
