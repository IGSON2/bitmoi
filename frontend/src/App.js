import { useState } from "react";
import { Route, BrowserRouter, Routes, Link } from "react-router-dom";
import Home from "./RoutingPages/home/Home";
import Competition from "./RoutingPages/competition/Competition";
import Practice from "./RoutingPages/practice/Practice";
import MyPage from "./RoutingPages/mypage/MyPage";
import Rank from "./RoutingPages/Ranking/Rank";
import Mobile from "./component/Mobile/Mobile";
import Wallet from "./component/Wallet/Wallet";
import SignUp from "./RoutingPages/SignUp/SignUp";
import AddBidding from "./RoutingPages/AdBidding/AdBidding";
import Login from "./RoutingPages/Login/Login";
import EmailLocate from "./RoutingPages/SignUp/EmailLocate";

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
          <Route path="/ad_bidding" element={<AddBidding />} />
          <Route path="/mypage" element={<MyPage />} />
          <Route path="/rank?page=1" element={<Rank />} />
          <Route path="/login" element={<Login />} />
          <Route path="/signup" element={<SignUp />} />
          <Route path="/goto/:domain" element={<EmailLocate />} />
          <Route path="/wallettest" element={<Wallet />} />
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
