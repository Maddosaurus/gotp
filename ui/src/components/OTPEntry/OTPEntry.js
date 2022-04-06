import React, { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";


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

    return (
      <>
        <h1>{entry.name}</h1>
        <small>Updated At: {entry.updateTime}</small><br />
        <small>Type: {entry.type}</small>
        <h3>123456</h3>
        <Link to={"/"}>Back to home page</Link>
      </>
    )
  }
