import React from "react";

const VerticalLine = () => {
  return (
    <div
      style={{
        height: "60%",
        width: "1%",
        textAlign: "left",
        borderLeft: "1px solid #aaa",
        lineHeight: "0.1em",
      }}
    >
      <span style={{ background: "#fff" }}></span>
    </div>
  );
};

export default VerticalLine;
