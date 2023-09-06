import React, { useState, useEffect } from "react";
import styles from "./Countdown.module.css";
import VerticalLine from "../lines/VerticalLine";
import Countbox from "./Countbox/Countbox";

function Countdown({ nextUnlock }) {
  const targetTime = new Date(nextUnlock);
  const [remainingSeconds, setRemainingSeconds] = useState(0);
  useEffect(() => {
    const interval = setInterval(updateCountdown, 1000);

    return () => {
      clearInterval(interval);
    };
  }, []);

  const updateCountdown = () => {
    const currentTime = new Date();
    const timeDiff = Math.max(targetTime - currentTime, 0);
    if (timeDiff === 0) {
      console.error(
        "target time is earlier than current time.",
        "target time :",
        targetTime,
        "current time",
        currentTime
      );
    }
    setRemainingSeconds(Math.floor(timeDiff / 1000));
  };

  const days = Math.floor(remainingSeconds / 86400);
  const hours = Math.floor((remainingSeconds % 86400) / 3600);
  const minutes = Math.floor((remainingSeconds % 3600) / 60);
  const seconds = remainingSeconds % 60;

  return (
    <div className={styles.wrapper}>
      <Countbox number={days.toString().padStart(2, "0")} unit={"DAYS"} />
      <VerticalLine />
      <Countbox number={hours.toString().padStart(2, "0")} unit={"HOURS"} />
      <VerticalLine />
      <Countbox number={minutes.toString().padStart(2, "0")} unit={"MINUTES"} />
      <VerticalLine />
      <Countbox number={seconds.toString().padStart(2, "0")} unit={"SECONDS"} />
    </div>
  );
}

export default Countdown;
