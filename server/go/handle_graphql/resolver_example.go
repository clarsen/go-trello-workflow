package handle_graphql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

// type Resolver struct{}
//
// func (r *Resolver) Mutation() MutationResolver {
// 	return &mutationResolver{r}
// }
// func (r *Resolver) Query() QueryResolver {
// 	return &queryResolver{r}
// }
//
// type mutationResolver struct{ *Resolver }
//
// func (r *mutationResolver) PrepareWeeklyReview(ctx context.Context, year *int, week *int) (*GenerateResult, error) {
// 	panic("not implemented")
// }
// func (r *mutationResolver) FinishWeeklyReview(ctx context.Context, year *int, week *int) (*bool, error) {
// 	panic("not implemented")
// }
// func (r *mutationResolver) SetDueDate(ctx context.Context, taskID string, due time.Time) (*Task, error) {
// 	panic("not implemented")
// }
// func (r *mutationResolver) SetDone(ctx context.Context, taskID string, done bool, status *string, nextDue *time.Time) (*Task, error) {
// 	panic("not implemented")
// }
//
// type queryResolver struct{ *Resolver }
//
// func (r *queryResolver) Tasks(ctx context.Context, dueBefore *int, inBoardList *BoardListInput) ([]Task, error) {
// 	panic("not implemented")
// }
// func (r *queryResolver) WeeklyVisualization(ctx context.Context, year *int, week *int) (*string, error) {
// 	panic("not implemented")
// }
// func (r *queryResolver) MonthlyGoals(ctx context.Context) ([]*MonthlyGoal, error) {
// 	panic("not implemented")
// }
