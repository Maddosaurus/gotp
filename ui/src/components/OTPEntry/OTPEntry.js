import React, { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";

const jsotp = require('jsotp');


export default function OTPEntry() {
  let params = useParams();
  const [entry, setEntry] = useState("");
  const [hotp, setHotp] = useState("Press Button");
  const [totp, setTotp] = useState("00000");

  const clickLog = () => {
    var updated_entry = entry;
    updated_entry.counter++;
    // updated_entry.updateTime = new Date().toISOString(); // FIXME: Test this! :D
    setEntry(updated_entry);
    const hotp_srv = jsotp.HOTP(entry.secretToken);
    setHotp(hotp_srv.at(entry.counter));
    fetch("http://localhost:8081/v1/otp/entries/" + params.uuid, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(updated_entry),
    })
  };

  const totp_ticker = () => {
    var totp_srv = jsotp.TOTP(entry.secretToken);
    var new_val = totp_srv.now();
    //console.log(totp +" !== "+ new_val + "? -> " + (totp !== new_val));
    //FIXME: Debug this!
    if (totp !== new_val) {
      setTotp(new_val);
    }
  }

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
      // const totp = jsotp.TOTP(entry.secretToken);
      // const n = totp.now();
      const interval = setInterval(() => {
        totp_ticker();
      }, 1000);
      return (
        <>
          <h1>{entry.name}</h1>
          <small>Updated At: {entry.updateTime}</small><br />
          <small>Type: {entry.type}</small>

          <h3>{totp}</h3>
          <h2>FIXME: Implement countdown</h2>
          <Link to={"/"}>Back to home page</Link>
        </>
      )
    } else if (entry.type === "HOTP"){
      return (
        <>
          <h1>{entry.name}</h1>
          <small>Updated At: {entry.updateTime}</small><br />
          <small>Type: {entry.type}</small>

          {/* <h3>{n}</h3> */}
          <h3>{hotp}</h3>
          <button onClick={clickLog}>Click Me</button>
          <br/><br/><br/>
          <Link to={"/"}>Back to home page</Link>
        </>
      )
    }
  }
