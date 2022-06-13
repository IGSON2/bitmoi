import { useState } from "react";
import { Route, BrowserRouter, Routes, Link } from "react-router-dom";
import Home from "./RoutingPages/home/Home";
import Competition from "./RoutingPages/competition/Competition";
import Practice from "./RoutingPages/practice/Practice";
import MyScore from "./RoutingPages/myscore/MyScore";
import Ranking from "./RoutingPages/Ranking/Ranking";
import Mobile from "./component/Mobile/Mobile";
import app from "./component/backendConn/fbase";

function App() {
  const [width, setWidth] = useState(window.innerWidth);
  window.onresize = (e) => {
    setWidth(e.target.innerWidth);
  };
  return (
    <div className="App" style={{ width: "100%", height: "100%" }}>
      {width < 700 ? <Mobile /> : null}
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/competition" element={<Competition />} />
          <Route path="/practice" element={<Practice />} />
          <Route path="/community" element={<div>This is board page.</div>} />
          <Route
            path="/ad_bidding"
            element={
              <h1
                style={{
                  width: "100%",
                  height: "100%",
                  textAlign: "center",
                  marginTop: "20%",
                }}
              >
                금방 준비해서 돌아오겠습니다!
              </h1>
            }
          />
          <Route path="/myscore" element={<MyScore />} />
          <Route path="/ranking" element={<Ranking />} />
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
