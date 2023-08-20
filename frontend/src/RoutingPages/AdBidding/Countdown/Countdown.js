import React, { useState, useEffect } from "react";
import styles from "./Countdown.module.css";
import VerticalLine from "../../../component/lines/VerticalLine";
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
    setRemainingSeconds(Math.floor(timeDiff / 1000));
  };

  const days = Math.floor(remainingSeconds / 86400);
  const hours = Math.floor(remainingSeconds / 3600);
  const minutes = Math.floor((remainingSeconds % 3600) / 60);
  const seconds = remainingSeconds % 60;

  return (
    <div className={styles.wrapper}>
      <Countbox number={days} unit={"DAYS"} />
      <VerticalLine />
      <Countbox number={hours} unit={"HOURS"} />
      <VerticalLine />
      <Countbox number={minutes} unit={"MINUTES"} />
      <VerticalLine />
      <Countbox number={seconds} unit={"SECONDS"} />
    </div>
  );
}

export default Countdown;
