import React, { useState } from "react"
import { Badge, Button, Collapse, Row, Col, Progress } from "reactstrap"
import moment from "moment"
import {
  FaTimes,
  FaThumbsUp,
  FaLink,
  FaExternalLinkAlt,
  FaFileAlt,
  FaListUl,
  FaForward,
} from "react-icons/fa"
import { GoTasklist } from "react-icons/go"

let iconSize = 20

function TaskReviewItem({ task, setDone, addComment, moveTaskToList }) {
  const [expanded, setExpanded] = useState(false)

  let color = ""
  let value = 0
  if (task.due) {
    let delta_days = moment().diff(moment.unix(task.due)) / (86400 * 1000)
    // console.log('task', task.title, 'task.due', moment.unix(task.due), 'delta_days from now', delta_days)
    if (delta_days < -3) {
      color = "info"
      value = ((100 * (14 + delta_days)) / 14).toFixed(0)
      if (value < 0) {
        value = 0
      }
    } else if (delta_days >= -3 && delta_days < 0) {
      color = "warning"
      value = ((100 * (14 + delta_days)) / 14).toFixed(0)
    } else if (delta_days >= 0) {
      color = "danger"
      value = ((delta_days / 7) * 100).toFixed(0)
      if (value > 100) {
        value = 100
      }
    }
  }
  let activitySinceCreate =
    moment.unix(task.dateLastActivity).diff(moment.unix(task.createdDate)) /
    (86400 * 1000)

  return (
    <React.Fragment key={task.id}>
      <Row key={"row0" + task.id}>
        <Col xs={12} lg={12} key={"2" + task.id}>
          <div>
            <div className="task">
              <Badge color="info">
                {moment.unix(task.dateLastActivity).format("YYYY-MM-DD")}
              </Badge>
              {activitySinceCreate > 30 && (
                <Badge color="dark">
                  {moment.unix(task.createdDate).format("YYYY-MM-DD")}
                </Badge>
              )}{" "}
              <a target="_blank" rel="noopener noreferrer" href={task.url}>
                <FaLink size={iconSize} />
              </a>
              <a
                target="_blank"
                rel="noopener noreferrer"
                href={`trello://x-callback-url/showCard?x-source=go-trello-workflow&id=${task.id}`}
              >
                <FaExternalLinkAlt size={iconSize} />
              </a>
              {task.desc != "" && (
                <FaFileAlt
                  size={iconSize}
                  onClick={() => {
                    setExpanded(!expanded)
                  }}
                />
              )}
              {task.checklistItems.length > 0 && (
                <FaListUl
                  size={iconSize}
                  onClick={() => {
                    setExpanded(!expanded)
                  }}
                />
              )}
              {task.title}
              <Badge color="secondary">{task.list.board}</Badge>
              {"/"}
              <Badge color="secondary">{task.list.list}</Badge>{" "}
              <FaTimes
                size={iconSize}
                onClick={() => {
                  setDone.mutation({
                    variables: {
                      taskId: task.id,
                      done: true,
                      titleComment: "won't do: ",
                    },
                    optimisticResponse: {
                      setDone: {
                        __typename: "Task",
                        id: task.id,
                        list: {
                          __typename: "BoardList",
                          board: "Kanban daily/weekly",
                          list: "Done this week",
                        },
                      },
                    },
                  })
                }}
              />{" "}
              <FaForward
                size={iconSize}
                onClick={() => {
                  console.log("addComment", addComment, "for", task)
                  addComment.mutation({
                    variables: {
                      taskId: task.id,
                      comment: "review later",
                    },
                    optimisticResponse: {
                      addComment: {
                        ...task,
                        dateLastActivity: moment().unix(),
                      },
                    },
                  })
                }}
              />{" "}
              <FaThumbsUp
                size={iconSize}
                onClick={() => {
                  setDone.mutation({
                    variables: {
                      taskId: task.id,
                      done: true,
                      titleComment: "Enough for now: ",
                    },
                    optimisticResponse: {
                      setDone: {
                        __typename: "Task",
                        id: task.id,
                        list: {
                          __typename: "BoardList",
                          board: "Kanban daily/weekly",
                          list: "Done this week",
                        },
                      },
                    },
                  })
                }}
              />{" "}
              <GoTasklist
                size={iconSize}
                onClick={() => {
                  moveTaskToList.mutation({
                    variables: {
                      taskID: task.id,
                      list: {
                        board: "Backlog (Personal)",
                        list: "Backlog",
                      },
                    },
                    optimisticResponse: {
                      moveTaskToList: {
                        ...task,
                        list: {
                          __typename: "BoardList",
                          board: "Backlog (Personal)",
                          list: "Backlog",
                        },
                      },
                    },
                  })
                }}
              />
            </div>
            <style jsx global>{`
              .periodicProgress {
                float: right;
                background-color: #888;
              }
            `}</style>
          </div>
        </Col>
      </Row>
      {expanded && (
        <React.Fragment>
          {task.checklistItems.length > 0 && (
            <Row key={"checkitems-" + task.id}>
              <Col xs={12} lg={12} key={"checkitems-" + task.id}>
                Checklist:
                <ul>
                  {task.checklistItems.map((item) => (
                    <li>{item}</li>
                  ))}
                </ul>
              </Col>
            </Row>
          )}
          {task.desc != "" && (
            <Row key={"description-" + task.id}>
              <Col xs={12} lg={12} key={"desc-" + task.id}>
                Description:
                <pre style={{ color: "white" }}>{task.desc}</pre>
              </Col>
            </Row>
          )}
        </React.Fragment>
      )}
    </React.Fragment>
  )
}

export default TaskReviewItem
