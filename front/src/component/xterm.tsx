import React, { useEffect } from "react";
import "xterm/css/xterm.css";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import Base64 from "crypto-js/enc-base64";
import Utf8 from "crypto-js/enc-utf8";

const msgInput = '1';
const msgResizeTerminal = '3';

interface Props {
    id: string;
}

export const XTermComponent = (props: Props) => {

    let { id } = props;
    const terminal = new Terminal();
    const fitAddon = new FitAddon();
    terminal.loadAddon(fitAddon);

    useEffect(() => {
        let terminalContainer = document.getElementById(id);
        const webSocket = new WebSocket(`ws://${window.location.host}/ws`);

        webSocket.onmessage = (event) => {
            terminal.write(Base64.parse(event.data).toString(Utf8));
        };

        webSocket.onopen = () => {
            terminal.open(terminalContainer);
            fitAddon.fit();
            terminal.write("welcome to bifrost🌈\r\n");
            terminal.focus();
        };

        webSocket.onclose = () => {
            terminal.write("\r\nwebTTY quit!");
        };

        webSocket.onerror = (event) => {
            // eslint-disable-next-line no-console
            console.error(event);
            webSocket.close();
        };

        terminal.onKey((event) => {
            webSocket.send(msgInput + Base64.stringify(Utf8.parse(event.key)));
        });

        terminal.onResize(({ cols, rows }) => {
            webSocket.send(msgResizeTerminal +
                Base64.stringify(
                    Utf8.parse(
                        JSON.stringify({
                            columns: cols,
                            rows: rows
                        })
                    )
                )
            );
        });
    })

    return <div id={id}></div>
}
