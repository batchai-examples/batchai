import { CommandLog } from "@/lib";
import React, { useEffect, useRef } from "react";
import ReactAnsi from "react-ansi";

interface LogViewerProps {
    logs: CommandLog[];
}

export function ExecutionLogViewer({ logs }: LogViewerProps) {
    const scrollRef = useRef<HTMLDivElement | null>(null);

    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
        }
    }, [logs.length]);

    return (
        <div ref={scrollRef} style={{ height: "160vh", overflowY: "auto", backgroundColor: "#2c2f3a" }}>
            <ReactAnsi autoScroll={false} log={logs.map((log) => log.message)} logStyle={{ fontSize: 13, backgroundColor: "#2c2f3a" }} />
        </div>
    );
};


export function AuditLogViewer({ logs }: LogViewerProps) {
    const scrollRef = useRef<HTMLDivElement | null>(null);

    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
        }
    }, [logs]);

    return (
        <div ref={scrollRef} style={{ height: "200vh", overflowY: "auto", backgroundColor: "#2c2f3a" }}>
            <ReactAnsi autoScroll={false} log={logs.map(log => `${log.timestamp}    ${log.message}`)} logStyle={{ fontSize: 13, backgroundColor: "#2c2f3a" }} />
        </div>
    );
};
