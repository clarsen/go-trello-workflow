package handle_graphql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

// type Resolver struct{}
//
// func (r *Resolver) MonthlyGoal() MonthlyGoalResolver {
// 	return &monthlyGoalResolver{r}
// }
// func (r *Resolver) Mutation() MutationResolver {
// 	return &mutationResolver{r}
// }
// func (r *Resolver) Query() QueryResolver {
// 	return &queryResolver{r}
// }
//
// type monthlyGoalResolver struct{ *Resolver }
//
// func (r *monthlyGoalResolver) WeeklyGoals(ctx context.Context, obj *MonthlyGoal) ([]WeeklyGoal, error) {
// 	panic("not implemented")
// }
//
// type mutationResolver struct{ *Resolver }
//
// func (r *mutationResolver) PrepareWeeklyReview(ctx context.Context, year *int, week *int) (*GenerateResult, error) {
// 	panic("not implemented")
// }
// func (r *mutationResolver) FinishWeeklyReview(ctx context.Context, year *int, week *int) (*FinishResult, error) {
// 	panic("not implemented")
// }
// func (r *mutationResolver) SetDueDate(ctx context.Context, taskID string, due time.Time) (*Task, error) {
// 	panic("not implemented")
// }
// func (r *mutationResolver) SetDone(ctx context.Context, taskID string, done bool, status *string, nextDue *time.Time) (*Task, error) {
// 	panic("not implemented")
// }
// func (r *mutationResolver) MoveTaskToList(ctx context.Context, taskID string, list BoardListInput) (*Task, error) {
// 	panic("not implemented")
// }
// func (r *mutationResolver) StartTimer(ctx context.Context, taskID string, checkitemID *string) (*Timer, error) {
// 	panic("not implemented")
// }
// func (r *mutationResolver) StopTimer(ctx context.Context, timerID string) (*bool, error) {
// 	panic("not implemented")
// }
// func (r *mutationResolver) SetGoalDone(ctx context.Context, taskID string, checkitemID string, done bool, status *string) ([]MonthlyGoal, error) {
// 	panic("not implemented")
// }
// func (r *mutationResolver) AddTask(ctx context.Context, title string, board *string, list *string) (*Task, error) {
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
// func (r *queryResolver) MonthlyGoals(ctx context.Context) ([]MonthlyGoal, error) {
// 	panic("not implemented")
// }
// func (r *queryResolver) ActiveTimer(ctx context.Context) (*Timer, error) {
// 	panic("not implemented")
// }
