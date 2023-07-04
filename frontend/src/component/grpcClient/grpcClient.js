import { BitmoiClient } from "./pb/service_bitmoi_grpc_web_pb";
import { CandlesRequest } from "./pb/candles_pb";

function grpcCleint() {
  const client = new BitmoiClient("bitmoi.co.kr:6000");
  const req = new CandlesRequest();
  const res = client.requestCandles(req);
  console.log(res);
}

export default grpcCleint;
