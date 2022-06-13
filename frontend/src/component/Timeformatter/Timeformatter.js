function Timeformatter(unixTime, isChart) {
  const scoreTime = new Date(parseInt(unixTime));
  var timeGap = 9;
  if (isChart) {
    timeGap = 0;
  }
  const formattedTime =
    // scoreTime.getUTCFullYear().toString().slice(2, 4) +
    // "." +
    (scoreTime.getUTCMonth() + 1).toString().padStart(2, "0") +
    "/" +
    scoreTime.getUTCDate().toString().padStart(2, "0") +
    " " +
    ((scoreTime.getUTCHours() + timeGap) % 24).toString().padStart(2, "0") +
    ":" +
    scoreTime.getUTCMinutes().toString().padStart(2, "0");
  return formattedTime;
}

export default Timeformatter;
