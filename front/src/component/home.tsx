import React from "react";
import "./home.css";
import { XTermComponent } from "./xterm";

export const HomeComponent = () => {

    return <div className="container">
        <XTermComponent id="xterm" />
    </div>
}