"use client";

import React, { useEffect, useState } from 'react';
import ReactAnsi from 'react-ansi'; 
import Box from '@mui/material/Box';
import Stepper from '@mui/material/Stepper';
import Step from '@mui/material/Step';
import StepLabel from '@mui/material/StepLabel';
import { CommandDetail, CommandRunStatus, CommandStatus, SessionState, useSession } from "@/lib";
import * as commandApi from '@/api/command.api';
import { UIContextType, useUIContext } from '@/lib/ui.context';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import Drawer from '@mui/material/Drawer';
import IconButton from '@mui/material/IconButton';
import Divider from '@mui/material/Divider';
import RefreshIcon from '@mui/icons-material/RefreshOutlined';
import Link from '@mui/material/Link';
import Toolbar from '@mui/material/Toolbar';
import StopIcon from '@mui/icons-material/StopOutlined';
import ResumeIcon from '@mui/icons-material/NavigateNextOutlined';
import ToolbarIcon from '@/components/toolbar.button';
import ResetIcon from '@mui/icons-material/UndoOutlined';
import DeleteIcon from '@mui/icons-material/DeleteOutlined';

interface Step {
  status: CommandRunStatus;
  label: string;
  needChanges: boolean;
}

const steps: Step[] = [
  {
    status: CommandRunStatus.Begin,
    label: "Begin",
    needChanges: false
  }, {
    status: CommandRunStatus.CheckedRemote,
    label: "Check the remote repository",
    needChanges: false
  }, {
    status: CommandRunStatus.Forked,
    label: "Fork it",
    needChanges: false
  }, {
    status: CommandRunStatus.ClonedOrPulled,
    label: "`git clone/pull`",
    needChanges: false
  }, {
    status: CommandRunStatus.CheckedOut,
    label: "`git checkout -b feature/batchai`",
    needChanges: false
  }, {
    status: CommandRunStatus.BatchAIExecuted,
    label: "Execute batchai",
    needChanges: false
  }, {
    status: CommandRunStatus.ChangesAdded,
    label: "`git add .`",
    needChanges: false
  }, {
    status: CommandRunStatus.ChangesCommited,
    label: "`git commit -m ...`",
    needChanges: false
  }, {
    status: CommandRunStatus.ChangesPushed,
    label: "`git push`",
    needChanges: true
  }, {
    status: CommandRunStatus.GetCommitId,
    label: "Get commit id",
    needChanges: true
  }, {
    status: CommandRunStatus.End,
    label: "End",
    needChanges: false
  },
];

async function refreshPage(s: SessionState, ui: UIContextType, id: number, 
          setCommand:React.Dispatch<React.SetStateAction<CommandDetail>>,
          setLog: React.Dispatch<React.SetStateAction<string>>,
          setActiveStep: React.Dispatch<React.SetStateAction<number>>): Promise<void> {
  ui.setLoading(true);
  try {
    const c = await commandApi.LoadCommand(s, ui, id);
    refreshCommand(s, ui, c, setCommand, setLog, setActiveStep);
  } catch (err) {
    ui.setError(err);
  } finally {
    ui.setLoading(false);
  };
}

async function refreshCommand(s: SessionState, ui: UIContextType, c: CommandDetail, 
          setCommand:React.Dispatch<React.SetStateAction<CommandDetail>>,
          setLog: React.Dispatch<React.SetStateAction<string>>,
          setActiveStep: React.Dispatch<React.SetStateAction<number>>): Promise<void> {
  ui.setLoading(true);
  try {
    setActiveStep(steps.findIndex(x => x.status === c.runStatus));
    setCommand(c);

    const log = await commandApi.LoadCommandLog(s, ui, c.id);
    //setLog(log.replace(/\n/g, '<br/>') || "...");
    setLog(log);
  } catch (err) {
    ui.setError(err);
  } finally {
    ui.setLoading(false);
  };
}

export default function CommandHome({ params }) {
  const [command, setCommand] = useState<CommandDetail>(null);
  const status = command?.status || '';
  const [log, setLog] = useState<string>("...");
  const repo = command?.repo;
  const owner = repo?.owner;
  const id = params.id;
  const s = useSession().state;
  const [activeStep, setActiveStep] = useState(0);
  const ui = useUIContext();

  const [showProgress, setShowProgress] = React.useState(false);
  const toggleProgress = (show: boolean) => () => {
    setShowProgress(show);
  };

  const handleCopy = () => {
    navigator?.clipboard?.writeText(command?.commandLine()).then(
      () => {
        alert('Copied to clipboard!');//TODO
      },
      (err) => {
        console.error('Failed to copy: ', err);
      }
    );
  };

  useEffect(() => {
    refreshPage(s, ui, id, setCommand, setLog, setActiveStep);
  }, [s, ui, id]);

  const onRefresh = async() => {
    refreshPage(s, ui, id, setCommand, setLog, setActiveStep);
  };

  const onReset = async() => {
    const c = await commandApi.ResetCommand(s, ui, id);
    refreshCommand(s, ui, c, setCommand, setLog, setActiveStep);
  };

  const onResume = async() => {
    const c = await commandApi.ResumeCommand(s, ui, id);
    refreshCommand(s, ui, c, setCommand, setLog, setActiveStep);
  };

  const onStop = async() => {
    const c = await commandApi.StopCommand(s, ui, id);
    refreshCommand(s, ui, c, setCommand, setLog, setActiveStep);
  };

  const onDelete = async() => {
    const c = await commandApi.RemoveCommand(s, ui, id);
    refreshCommand(s, ui, c, setCommand, setLog, setActiveStep);
  };

  return (
    <>
      <Box sx={{ mb: 2 }}>
        <Typography variant="body2">{command?.command}</Typography><Typography variant="h5" component="a" href={repo?.repoUrl}>{owner?.name} / {repo?.name}</Typography>
        <Typography variant="body2">{status}</Typography>
        <Button onClick={toggleProgress(true)}>Detailed progress</Button>

        <Box sx={{ mt: 1, display: 'flex', flexDirection: 'row' }}>
          {command?.hasChanges ?
            <>
              <Typography variant="body2">Changes:</Typography>
              <Link sx={{ ml: 2 }} href={command?.commitId ? command.commitUrl : '#'}>{command?.commitId ? command.commitUrl : 'N/A'}</Link>
            </>
          :
            <Typography variant="body2">No changes</Typography>
          }
        </Box>
      </Box>

      <Divider sx={{ mb: 2 }} />

      <Toolbar sx={{ mb: 2 }}>
        <ToolbarIcon label='Refresh' enabled={true} onClick={onRefresh}>
          <RefreshIcon/>
        </ToolbarIcon>
        <ToolbarIcon label='Reset' enabled={status!==CommandStatus.Running} onClick={onReset}>
          <ResetIcon/>
        </ToolbarIcon>
        <ToolbarIcon label='Stop' enabled={status===CommandStatus.Running} onClick={onStop}>
          <StopIcon/>
        </ToolbarIcon>
        <ToolbarIcon label='Resume' enabled={status===CommandStatus.Pending || status===CommandStatus.Failed} onClick={onResume}>
          <ResumeIcon/>
        </ToolbarIcon>
        <ToolbarIcon label='Delete' enabled={status!==CommandStatus.Running} onClick={onDelete}>
          <DeleteIcon/>
        </ToolbarIcon>
      </Toolbar>

      <Box sx={{ mb: 2 }}>
        <Typography sx={{ mt: 2 }} variant="body2">Command Line & Execution Log</Typography>
        <div
          style={{
            paddingTop: '5px', paddingBottom: '5px', paddingLeft: '12px', paddingRight: '12px',
            width: '100%', background: '#404040', color: 'white', display: 'flex', alignItems: 'center',
            justifyContent: 'space-between',
          }}
        >
          <div style={{ paddingLeft: 50 }}>
            $ <span style={{ color: '#80ff80' }}>batchai</span>
            &nbsp;&nbsp;{command?.globalOptions.join(' ')}
            <span style={{ color: '#80ff80' }}>&nbsp;&nbsp;{command?.command}</span>
            &nbsp;&nbsp;{command?.commandOptions}
            &nbsp;&nbsp;.
            &nbsp;&nbsp;{command?.targetPaths.join(' ')}
          </div>
          <IconButton onClick={handleCopy} color="primary" aria-label="copy command">
            <ContentCopyIcon />
          </IconButton>
        </div>
          <ReactAnsi log={log} logStyle={{fontSize: 12}}/>
      </Box>

      <Drawer anchor='right' open={showProgress} onClose={toggleProgress(false)}      >
        <Stepper sx={{ margin: 4 }} activeStep={activeStep} orientation="vertical">
          {steps.map((step, index) => {
            return (
              <Step key={step.status} completed={index <= activeStep && (!step.needChanges || (step.needChanges && command?.hasChanges))}>
                <StepLabel error={index === (activeStep + 1) && status==CommandStatus.Failed}>{step.label}</StepLabel>
              </Step>
            );
          })}
        </Stepper>
      </Drawer>
    </>
  );
}
