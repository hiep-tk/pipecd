import {
  Box,
  Collapse,
  IconButton,
  List,
  ListItem,
  makeStyles,
  Typography,
} from "@material-ui/core";
import dayjs from "dayjs";
import React, { FC, useMemo, useState } from "react";
import { ListDeploymentTracesResponse } from "~~/api_client/service_pb";
import MoreHorizIcon from "@material-ui/icons/MoreHoriz";
import { ArrowDropDown } from "@material-ui/icons";
import DeploymentItem from "./deployment-item";
import { Link as RouterLink } from "react-router-dom";
import { PAGE_PATH_DEPLOYMENTS } from "~/constants/path";
import clsx from "clsx";

const useStyles = makeStyles((theme) => ({
  title: {
    display: "inline",
  },
  btnMore: {
    display: "inline-flex",
    padding: "0 1px",
    borderRadius: 5,
    marginLeft: 5,
    marginBottom: 4,
  },
  btnMoreActive: {
    backgroundColor: theme.palette.grey[300],
  },
  btnRotate: {
    transform: "rotate(180deg)",
  },
  traceItem: {
    padding: theme.spacing(2),
    paddingRight: theme.spacing(0),
    borderBottom: `1px solid ${theme.palette.grey[300]}`,
    backgroundColor: theme.palette.background.paper,
    "&:hover": {
      backgroundColor: theme.palette.grey[100],
    },
  },
  traceStickyTop: {
    position: "sticky",
    top: 0,
    zIndex: 50,
    paddingBottom: theme.spacing(1),
  },
  commitMessageWrap: {
    maxHeight: "20svh",
    overflow: "hidden auto",
    borderLeft: `0.5px  solid ${theme.palette.grey[500]}`,
    paddingLeft: theme.spacing(1),
    paddingTop: theme.spacing(1),
    marginBottom: theme.spacing(1),
    marginLeft: theme.spacing(1),
  },
  commitMessage: {
    whiteSpace: "pre-wrap",
  },
  emptyPlaceholder: {
    padding: theme.spacing(2),
    border: `1px solid ${theme.palette.grey[300]}`,
    borderTop: "none",
    backgroundColor: theme.palette.background.default,
  },
}));

type Props = {
  trace: ListDeploymentTracesResponse.DeploymentTraceRes.AsObject["trace"];
  deploymentList: ListDeploymentTracesResponse.DeploymentTraceRes.AsObject["deploymentsList"];
};

const DeploymentTraceItem: FC<Props> = ({ trace, deploymentList }) => {
  const classes = useStyles();
  const [visibleMessage, setVisibleMessage] = useState(false);
  const [visibleDeployments, setVisibleDeployments] = useState(false);

  const onViewCommitMessage = (
    e: React.MouseEvent<HTMLButtonElement>
  ): void => {
    e.stopPropagation();
    setVisibleMessage(!visibleMessage);
  };

  const timeStampCommit = useMemo(() => {
    if (!trace?.commitTimestamp) return "-";
    const timeStamp = trace.commitTimestamp * 1000;
    const diff = dayjs().diff(timeStamp, "month");
    const date = dayjs(timeStamp);
    const isCurrentYear = dayjs().isSame(date, "year");

    if (!isCurrentYear) {
      return date.format("MMM D, YYYY");
    }
    if (diff > 1) {
      return date.format("MMM D");
    }

    return date.fromNow();
  }, [trace?.commitTimestamp]);

  return (
    <Box flex={1} width={"100%"}>
      <Box
        className={clsx(classes.traceItem, {
          [classes.traceStickyTop]: visibleDeployments,
        })}
      >
        <Box
          display="flex"
          gridColumnGap={10}
          alignItems={"start"}
          justifyContent={"space-between"}
          pr={1}
        >
          <Box overflow={"hidden"} flex={1}>
            <Box>
              <Typography variant="h6" className={classes.title}>
                {trace?.title || `Title of commit ${trace?.commitHash}`}
              </Typography>
              {trace?.commitMessage && (
                <IconButton
                  size="small"
                  aria-label="btn-commit-message"
                  className={clsx(classes.btnMore, {
                    [classes.btnMoreActive]: visibleMessage,
                  })}
                  onClick={onViewCommitMessage}
                  title={
                    visibleMessage
                      ? "Hide commit message"
                      : "View commit message"
                  }
                >
                  <MoreHorizIcon />
                </IconButton>
              )}
            </Box>

            <Box display="flex">
              <RouterLink to={trace?.commitUrl || "#"} target="_blank">
                <Typography variant="body2" color="textSecondary">
                  {trace?.commitHash}
                </Typography>
              </RouterLink>
            </Box>
          </Box>

          <IconButton
            aria-label="expand"
            className={visibleDeployments ? classes.btnRotate : ""}
            onClick={() => setVisibleDeployments(!visibleDeployments)}
          >
            <ArrowDropDown />
          </IconButton>
        </Box>

        {visibleMessage && (
          <Box className={classes.commitMessageWrap}>
            <Typography
              variant="body2"
              color="textSecondary"
              className={classes.commitMessage}
            >
              {trace?.commitMessage}
            </Typography>
          </Box>
        )}

        <Box display={"flex"} gridColumnGap={3}>
          {trace?.author && (
            <Typography variant="body2" color="textSecondary">
              {trace?.author} authored
            </Typography>
          )}
          <Typography
            variant="body2"
            color="textSecondary"
            title={dayjs(trace?.commitTimestamp).format("MMM D, YYYY h:mm A")}
          >
            {timeStampCommit}
          </Typography>
        </Box>
      </Box>

      <Collapse in={visibleDeployments} unmountOnExit key={trace?.id}>
        {deploymentList.length === 0 && (
          <Box className={classes.emptyPlaceholder}>
            <Typography variant="body2" color="textSecondary" align="center">
              No deployment triggered
            </Typography>
          </Box>
        )}
        <List>
          {deploymentList.map((deployment) => (
            <ListItem
              key={deployment?.id}
              button
              dense
              divider
              component={RouterLink}
              to={`${PAGE_PATH_DEPLOYMENTS}/${deployment.id}`}
            >
              <DeploymentItem key={deployment.id} deployment={deployment} />
            </ListItem>
          ))}
        </List>
      </Collapse>
    </Box>
  );
};

export default DeploymentTraceItem;
