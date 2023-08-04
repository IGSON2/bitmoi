import axiosClient from "../../component/backendConn/axiosClient";

function Freetoken() {
  const getFreeToken = async () => {
    await axiosClient.post("/freetoken", {
      addr: "0x6655992CEDa8A8Faf070208A96a1051144D77E3D",
    });
  };

  return (
    <div>
      <button onClick={getFreeToken}>토큰 발급받기</button>
    </div>
  );
}

export default Freetoken;
