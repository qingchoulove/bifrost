import React, { useEffect } from "react";
import "./home.css";
import { XTermComponent } from "./xterm";

export const HomeComponent = () => {

    useEffect(() => {
        console.log("did mount");
        return function cleanup() {
            console.log("cleanup");
        }
    })

    return <div className="container">
        <XTermComponent id="xterm" />
    </div>
}