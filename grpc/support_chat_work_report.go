package grpc

import (
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_report"
	"github.com/tel4vn/fins-microservices/model"
)

func convertChatWorkReportToPbChatWorkReport(data *[]model.ChatWorkReport) (result []*pb.WorkReport) {
	result = make([]*pb.WorkReport, 0)
	for _, report := range *data {
		tmp := &pb.WorkReport{
			UserId:       report.UserId,
			UserFullname: report.UserFullname,
			Total:        int32(report.Total),
			Facebook: &pb.ChannelWorkPerformance{
				TotalChannels: int32(report.Facebook.TotalChannels),
				ReceivingTime: &pb.PerformanceMetrics{
					Fastest: int32(report.Facebook.ReceivingTime.Fastest),
					Average: int32(report.Facebook.ReceivingTime.Average),
					Slowest: int32(report.Facebook.ReceivingTime.Slowest),
				},
				ReplyingTime: &pb.PerformanceMetrics{
					Fastest: int32(report.Facebook.ReplyingTime.Fastest),
					Average: int32(report.Facebook.ReplyingTime.Average),
					Slowest: int32(report.Facebook.ReplyingTime.Slowest),
				},
			},
			Zalo: &pb.ChannelWorkPerformance{
				TotalChannels: int32(report.Zalo.TotalChannels),
				ReceivingTime: &pb.PerformanceMetrics{
					Fastest: int32(report.Zalo.ReceivingTime.Fastest),
					Average: int32(report.Zalo.ReceivingTime.Average),
					Slowest: int32(report.Zalo.ReceivingTime.Slowest),
				},
				ReplyingTime: &pb.PerformanceMetrics{
					Fastest: int32(report.Zalo.ReplyingTime.Fastest),
					Average: int32(report.Zalo.ReplyingTime.Average),
					Slowest: int32(report.Zalo.ReplyingTime.Slowest),
				},
			},
		}

		result = append(result, tmp)
	}
	return
}

func convertChatGeneralReportToPbChatGeneralReport(data *[]model.ChatGeneralReport) (result []*pb.GeneralReport) {
	result = make([]*pb.GeneralReport, 0)
	for _, report := range *data {
		tmp := &pb.GeneralReport{
			Channel:            report.Channel,
			OaName:             report.OaName,
			TotalConversations: int32(report.TotalConversations),
			Fresh: &pb.QuantityRatio{
				Quantity:   int32(report.Fresh.Quantity),
				Percentage: int32(report.Fresh.Percent),
			},
			Processing: &pb.QuantityRatio{
				Quantity:   int32(report.Processing.Quantity),
				Percentage: int32(report.Processing.Percent),
			},
			Resolved: &pb.QuantityRatio{
				Quantity:   int32(report.Resolved.Quantity),
				Percentage: int32(report.Resolved.Percent),
			},
		}
		result = append(result, tmp)
	}
	return
}
