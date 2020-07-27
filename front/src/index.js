import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import 'xterm/css/xterm.css';

const terminal = new Terminal();
const fitAddon = new FitAddon();
const webSocket = new WebSocket('ws://127.0.0.1:8080/ws');
const terminalContainer = document.getElementById('terminal');
terminal.loadAddon(fitAddon);

webSocket.onerror = (event) => {
    // eslint-disable-next-line no-console
    console.error(event);
    webSocket.close();
};

webSocket.onopen = () => {
    terminal.open(terminalContainer);
    fitAddon.fit();
    terminal.focus();
};

webSocket.onmessage = (event) => {
    terminal.write(window.atob(event.data));
};

webSocket.onclose = () => {
    terminal.write('\r\n');
    terminal.write('webTTY quit!');
};

terminal.onKey((event) => {
    webSocket.send(window.btoa(event.key));
});
