import React, { useState } from "react"
import {
  Row,
  Col,
  Navbar,
  Nav,
  NavItem,
  NavLink,
  TabContent,
  TabPane,
} from "reactstrap"
import { FaSync } from "react-icons/fa"
import TaskReviewList from "./TaskReviewList"
import classnames from "classnames"

export default function OldTasksReview({
  allRefetch,
  loadingAll,
  queryAllError,
  allTasks,
  setDone,
  addComment,
  moveTaskToList
}) {
  const [activeTab, setActiveTab] = useState("backlog")

  return (
    <React.Fragment>
      <Navbar expand="lg">
        <Nav className="mr-auto" tabs>
          <NavItem>
            <NavLink
              className={classnames({
                active: activeTab === "backlog",
              })}
              onClick={() => setActiveTab("backlog")}
            >
              Backlog
            </NavLink>
          </NavItem>
          <NavItem>
            <NavLink
              className={classnames({
                active: activeTab === "someday/maybe",
              })}
              onClick={() => setActiveTab("someday/maybe")}
            >
              Someday/Maybe
            </NavLink>
          </NavItem>
          <NavItem>
            <NavLink
              className={classnames({
                active: activeTab === "movies/tv",
              })}
              onClick={() => setActiveTab("movies/tv")}
            >
              Movies, TV
            </NavLink>
          </NavItem>
        </Nav>
      </Navbar>
      <FaSync
        size={25}
        onClick={() => {
          allRefetch()
        }}
      />{" "}
      <TabContent activeTab={activeTab}>
        <TabPane tabId="backlog">
          <Row>
            <Col lg={12}>
              <TaskReviewList
                loading={loadingAll}
                error={queryAllError}
                data={allTasks}
                listTitle={"Backlog (Personal)"}
                boardFilter={[
                  "Backlog (Personal)",
                ]}
                setDone={setDone}
                addComment={addComment}
                moveTaskToList={moveTaskToList}
              />{" "}
            </Col>{" "}
          </Row>{" "}
        </TabPane>
        <TabPane tabId="someday/maybe">
          <Row>
            <Col lg={12}>
              <TaskReviewList
                loading={loadingAll}
                error={queryAllError}
                data={allTasks}
                listTitle={"Someday/Maybe"}
                boardFilter={[
                  "Someday/Maybe",
                ]}
                setDone={setDone}
                addComment={addComment}
                moveTaskToList={moveTaskToList}
              />{" "}
            </Col>{" "}
          </Row>{" "}
        </TabPane>
        <TabPane tabId="movies/tv">
          <Row>
            <Col lg={12}>
              <TaskReviewList
                loading={loadingAll}
                error={queryAllError}
                data={allTasks}
                listTitle={"Movies, TV"}
                boardFilter={[
                  "Movies, TV",
                ]}
                setDone={setDone}
                addComment={addComment}
                moveTaskToList={moveTaskToList}
              />{" "}
            </Col>{" "}
          </Row>{" "}
        </TabPane>
      </TabContent>
    </React.Fragment>
  )
}
