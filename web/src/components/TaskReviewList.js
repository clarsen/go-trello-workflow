import React, { useState } from "react"
import TaskReviewItem from "./TaskReviewItem"
import { Container, Spinner } from "reactstrap"

export default function TaskReviewList({
  listTitle,
  loading,
  error,
  data,
  boardFilter,
  listFilter,
  setDone,
  addComment,
  moveTaskToList
}) {
  const [showControls, setShowControls] = useState(false)

  return (
    <React.Fragment>
      {listTitle && (
        <div
          className="listTitle"
          onClick={() => {
            setShowControls(!showControls)
          }}
        >
          {listTitle}
        </div>
      )}
      {loading && <Spinner color="primary" />}
      {error && <div>Tasks: {error.message}</div>}
      {!loading && !error && (
        <Container>
          {data.tasks
            .filter(
              (t) =>
                (!listFilter ||
                  listFilter.count == 0 ||
                  listFilter.includes(t.list.list)) &&
                (!boardFilter ||
                  boardFilter.count == 0 ||
                  boardFilter.includes(t.list.board))
            )
            .sort((a, b) => a.createdDate - b.createdDate)
            .sort((a, b) => a.dateLastActivity - b.dateLastActivity)
            .map((t) => (
              <TaskReviewItem
                key={t.id}
                setDone={setDone}
                addComment={addComment}
                moveTaskToList={moveTaskToList}
                task={t}
              />
            ))}
        </Container>
      )}
      <style global jsx>{`
        #newTaskTitle {
          width: 50%;
        }
        .listSubGroupTitle {
          background: #999;
          width: 100%;
        }
        .listTitle {
          background: #bbb;
          width: 100%;
          color: #fff;
        }
      `}</style>
    </React.Fragment>
  )
}
