import React, { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";

const jsotp = require('jsotp');


export default function OTPEntry() {
  let params = useParams();
  const [entry, setEntry] = useState("");//useState({"secretToken":"2222222222222222"});
  const [hotp, setHotp] = useState("Press Button");
  const [totp, setTotp] = useState("-");
  const [secs, setSecs] = useState("-");

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

  // see https://eight-bites.blog/en/2021/05/setinterval-setstate/
  // TODO / FIXME: Improve reload behaviour
  const totp_ticker = () => {
    var totp_srv = jsotp.TOTP(entry.secretToken);
    setInterval(() => {
      var new_val = totp_srv.now();
      setTotp(new_val);
      var currentSeconds = 30-parseInt(new Date().getTime()/1000%30);
      setSecs(currentSeconds);
    }, 1000)
  }

  useEffect(() => {
        fetch("http://localhost:8081/v1/otp/entries/" + params.uuid)
        .then (res => res.json())
        .then((data) => {
            setEntry(data);
        })
        .then(() => {
            // FIXME: This doesn't work. We would need to wait for the setEntry to happen,
            // as the totp_ticker tries do get the undefined secretToken right now.
            totp_ticker()
        })
        .catch(console.log)
        return () => clearInterval(totp_ticker)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [totp]) // This will lead to useEffect only being called once, vs an infinite loop
           // https://stackoverflow.com/questions/53715465/can-i-set-state-inside-a-useeffect-hook


    if (entry.type === "TOTP") {
      return (
        <>
          <h1>{entry.name}</h1>
          <small>Updated At: {entry.updateTime}</small><br />
          <small>Type: {entry.type}</small>

          <h3>{totp}<br/>[{secs}]</h3>
          <Link to={"/"}>Back to home page</Link>
        </>
      )
    } else if (entry.type === "HOTP"){
      return (
        <>
          <h1>{entry.name}</h1>
          <small>Updated At: {entry.updateTime}</small><br />
          <small>Type: {entry.type}</small>

          <h3>{hotp}</h3>
          <button onClick={clickLog}>Click Me</button>
          <br/><br/><br/>
          <Link to={"/"}>Back to home page</Link>
        </>
      )
    }
  }
