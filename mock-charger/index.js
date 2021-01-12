var http = require("http");

totalConsumption = 1;

function fakeData() {
  return {
    ChgState: Math.random() > 0.5 ? "Idle" : "Charging", // todo real names??
    Tariff: "T1",
    Price: 280,
    ChgDuration: 11811,
    ChgNrg: totalConsumption++, // 1000 * P [kWh]
    NrgDemand: Math.floor(4000 + Math.random() * 1000),
    Solar: 0,
    EmTime: 1440,
    RemTime: 1440,
    ActPwr: 0,
    ActCurr: Math.floor(9000 + Math.random() * 5000),
    MaxCurrT1: 16,
    BeginH_T1: 4,
    BeginM_T1: 30,
    PriceT1: 280,
    MaxCurrT2: 16,
    BeginH_T2: 22,
    BeginM_T2: 0,
    PriceT2: 200,
    RemoteCurr: 16,
    SolarPrice: 0,
    ExcessNrg: false,
    TMaxCurrT1: 16,
    TBeginH_T1: 4,
    TBeginM_T1: 30,
    TPriceT1: 280,
    TMaxCurrT2: 16,
    TBeginH_T2: 22,
    TBeginM_T2: 0,
    TPriceT2: 200,
    TRemoteCurr: 16,
    TSolarPrice: 0,
    TExcessNrg: true,
    HCCP: "A11",
  };
}

//create a server object:
http
  .createServer(function (req, res) {
    console.log(req);
    res.write(JSON.stringify(fakeData())); //write a response to the client
    res.end(); //end the response
  })
  .listen(25000);
