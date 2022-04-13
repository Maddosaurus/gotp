import React, { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";

const jsotp = require('jsotp');


export default function OTPEntry() {
  let params = useParams();
  const [entry, setEntry] = useState("");

  useEffect(() => {
        fetch("http://localhost:8081/v1/otp/entries/" + params.uuid)
        .then (res => res.json())
        .then((data) => {
            setEntry(data)
        })
        .catch(console.log)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []) // This will lead to useEffect only being called once, vs an infinite loop
           // https://stackoverflow.com/questions/53715465/can-i-set-state-inside-a-useeffect-hook


    if (entry.type === "TOTP") {
      const totp = jsotp.TOTP(entry.secretToken);
      const n = totp.now();
      return (
        <>
          <h1>{entry.name}</h1>
          <small>Updated At: {entry.updateTime}</small><br />
          <small>Type: {entry.type}</small>

          <h3>{n}</h3>
          <h2>FIXME: Implement dynamic update and countdown</h2>
          <Link to={"/"}>Back to home page</Link>
        </>
      )
    } else if (entry.type === "HOTP"){
      const hotp = jsotp.HOTP(entry.secretToken);
      const n = hotp.at(entry.counter);
      return (
        <>
          <h1>{entry.name}</h1>
          <small>Updated At: {entry.updateTime}</small><br />
          <small>Type: {entry.type}</small>

          <h3>{n}</h3>
          <h2>FIXME: Implement actual counter update to server :3</h2>
          <h2>Idea: Tie this (and text rendering) to a callback button function</h2>
          <Link to={"/"}>Back to home page</Link>
        </>
      )
    }
  }
