package validate

import (
	"github.com/SashaMelva/calendar_service/internal/storage"
)

// type Validator func(req interface{}) error

// func UnaryServerRequestValidatorInterceptor(validator Validator) grpc.UnaryServerInterceptor {
// 	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 		if err := validator(req); err != nil {
// 			return nil, status.Errorf(codes.InvalidArgument, "%s is rejected by validate middleware. Error: %v", info.FullMethod, err)
// 		}
// 		return handler(ctx, req)
// 	}
// }

// func ValidateReq(req interface{}) error {
// 	switch r := req.(type) {
// 	case *pb.SubmitVoteRequest:
// 		if r.Vote.Passport == "" || r.Vote.CandidateId == 0 {
// 			return errors.New("middleware validator: passport or candidate_id wrong")
// 		}
// 	}
// 	return nil
// }

func ValidEvent(event *storage.Event) (string, string) {
	err := ""

	if event.Title == "" {
		err += "row title empty;"
	}

	if event.DateTimeStart == nil {
		err += "row date start empty;"
	} else if storage.Date(event.DateTimeStart) == "" {
		err += "date start param empty date;"
	}

	if event.DateTimeEnd == nil {
		err += "row date end empty;"
	} else if storage.Date(event.DateTimeEnd) == "" {
		err += "date end param empty date;"
	}

	if err != "" {
		return "Err", err
	}

	return "OK", ""
}
