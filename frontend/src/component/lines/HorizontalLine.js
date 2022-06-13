import React from "react";

const HorizontalLine = () => {
  return (
    <div
      style={{
        width: "100%",
        height: "1%",
        textAlign: "left",
        borderBottom: "1px solid #aaa",
        lineHeight: "0.1em",
      }}
    >
      <span style={{ background: "#fff" }}></span>
    </div>
  );
};

export default HorizontalLine;
