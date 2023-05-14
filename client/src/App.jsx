import "./App.css";
import * as forge from "node-forge";
import React from "react";

const N =
  BigInt(971);
const G = 3;

const URL = "http://127.0.0.1:1323"

function App() {
  const [privatePrimeNumber, setPrivatePrimeNumber] = React.useState(0);
  const [publicKey, setPublicKey] = React.useState(0);

  const [allPublicKeys,setAllPublicKeys] = React.useState({})

  React.useEffect(() => {

    // Investigate how big of a number can we use here
    var bits = 20;
    var options = {
      algorithm: {
        name: "PRIMEINC",
        workers: -1, 
      },
    };
    forge.prime.generateProbablePrime(bits, options, function (err, num) {
      setPrivatePrimeNumber(num);
    });
  }, []);

  // Generate public key and send to server
  React.useEffect(() => {
    if (privatePrimeNumber) {
      const publicKey = BigInt(BigInt(G) ** BigInt(privatePrimeNumber) ) % N
      setPublicKey(publicKey);
    }
  }, [privatePrimeNumber]);

  React.useEffect(() => {
    if(publicKey){
    fetch(`${URL}/api/sync-public-key` , {
      method:'POST',
      body : JSON.stringify({
        publicKey: publicKey.toString(),
        name: "bob"
      }),
      headers: {
        'Content-Type' :'application/json'
      }
    }).then(data => data.json()).then(keys => {
      setAllPublicKeys({
        ...allPublicKeys,
        ...keys
      })
    })
  }
  }, [publicKey]);

  return (
    <>
      <div className="flex items-center flex-col justify-center h-screen w-screen w-full p-2 rounded-lg">
        <div className="h-2/3 text-center border h-full w-full">
          {" "}
          Chat Window
        </div>

        <div className="w-full">
          <input
            className="h-10  w-full border outline-0 mt-4 p-2 rounded-lg"
            placeholder="Enter text"
          />
          <button className="border p-2 rounded-lg mt-3">Send</button>
          <button className="border p-2 rounded-lg ml-3">Handshake</button>
        </div>
      </div>
    </>
  );
}

export default App;
