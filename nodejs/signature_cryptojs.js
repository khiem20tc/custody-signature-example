
// Replace these with your actual API keys and secrets.
const apiKey = "5WWCMPHwIEeB";
const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAu265+cwNqwGIIRR5Ub/XCQ05fxYGnCGxIRdMM7lYlOvFiaQu
vBbjdcPtIsQuKIsuLfgHPVwibqzoNSMO91V6bN7ArkhrnXOb2MXNtu+ZodQ17MA9
+5MW6+pQHHfycvx/Gk1ldLJgtJiUSFaVujE7O4q10gtcMKT4J+HnrCFlEloQtxlT
ElEafUyYk07hbbz0U4saO935/kAoT+PdMkAxTUD1xMxSXEhk7CociuS/gFSPxg8H
KSzHp6iXgozOtOkvl99LD3yZ87EjxhzMKTyFQGx/WjooNzK5yyruPyqTlrgrfaIy
1CZOlRXAP4Ttk2IoEahI9f0sLdXssxXZuwWWkwIDAQABAoIBAD7CPJNft9Pil2o8
KMMusRney7m57kypG14xJtrK3NZAe8wypVNldpQgHm7dsXbx42yQ+Bublgvo6Xeh
XYmDnZKGo423whDefPiAgvkWESMWo1e6pwZtoecsddaScyP9V7G+6JHCiI7v5/aw
x0Go6mRtdaP3Gc9P7aetBJ2mMOmLnAvZkZz+ePMTys3Y0PpGBs/fnyWFZUyuI+w4
Eouz3oX/RNrOfTAymo9GvZjL0pZ4N9UTu3XEGQYd97eUA1J/9HQ8m9JaGl454bxL
1iZmXlfuPnCEJUF58pqDbZc+eejKVy5G496OQsOUxq3e5KEEuJqBh/+rqtBKecr4
EvmiaIECgYEA2sGbqhtT/N2plviHGoL8s7ULnaQIoH4dE+t1l81IjhRUNhLOD8Kc
EB8RBSpKcT7iiE+Kaz0ibhjIQQusHvd0q4n9LE9odMV2s+R2kk/4YMtpjmA8wBDs
1lHl4PpVQx7fq976PidKNTH7yaUIB19bREBWhGltnlRBpivrHf2S18sCgYEA21fj
C1mK59kqQ5GXijeqTrx5hm0JNRQFxHnBoLgMVJ42QiTvU+GhwMKhS0X80jIDdGip
7ydPJaRCo01ldYU7sOV1q/fpz/fEqC4FNLCFsm31VZldX976tihImNtq/+8pybch
ccPjFVo9lqGEDm3ID5kEdmtLBaMzSGyjzetNk1kCgYEAxEnxmePHqyBjKjp7UEi0
47PSZnNn4ksHYHZpH/tt3T9UiOi6yd2AF98ocJAQGCmrL1DgDXXfzRajqeoFWgwF
Pl8lM3tVaWI+LxETbBoh7wjXAJBOMrF9MppuQT+e/glX/mqn9NlgdvcQzVEuMR9Z
T5bDizDm0akc9zR1VoXQG50CgYAovqSwYQvKka6mKo9p33lFcwFoFS0WrQd9PdjY
EBhKR7FwjAfhHxK7CeyIXRHfweaeYyreAAFVzrOKPkBQmlVCQP2g2kaWmUHws8vH
w9qyEHb4Vargujz8RXNm4at4q2apz9jolyjBuKekKZCsVXxKWXRYwwmGnJBULcon
4EPi0QKBgGSyxEosiG5rKjwYW5RP3Ag5invC6yZ/SRHyTYZ8BObjGyo80zS9rMmx
t0hAvBYb1D8CUdq7zW+9RbLxWPuGnXgoIah/6FqNKQ4liAYbICprA+762+zDxgyF
dZ9Mh43LQ0onvSii1GsYSaIh5E9R2X5H1UXT1b+M4Ix+YXLeSUCZ
-----END RSA PRIVATE KEY-----`;

const forge = require("node-forge");

const signData = (message, privateKeyPEM) => {
  try {
    const privateKey = forge.pki.privateKeyFromPem(privateKeyPEM);

    const md = forge.md.sha256.create();
    md.update(message, "utf8");

    const signature = privateKey.sign(md);
    const signatureBase64 = forge.util.encode64(signature);

    console.log("Signature:", signatureBase64);

    return signatureBase64;
  } catch (err) {
    console.log({ err });
    return "";
  }
};

const createSignature = (method, path, body) => {
  // Define the request body (if applicable).
  const requestBody = JSON.stringify(body) || "";
  // Create a timestamp for the request.
  const timestamp = Date.now().toString();

  // Generate a nonce (random string) for each request.
  const nonce = Math.random().toString(36).substring(2);

  // Create a string to sign based on your API requirements.
  const stringToSign = `${apiKey}${timestamp}${nonce}${path}${method}${requestBody}`;

  const signature = signData(stringToSign, privateKey);

  // Define request headers with API key, timestamp, nonce, and signature.
  const headers = {
    "X-Aegis-Api-Key": apiKey,
    "X-Aegis-Api-Timestamp": timestamp,
    "X-Aegis-Api-Nonce": nonce,
    "X-Aegis-Api-Signature": signature,
  };
  return headers;
};

const test = async () => {
  const postUrl = "http://127.0.0.1:3000/api/v1/address/deposit-address?depositLabel=poolETH1a";
  const postBody = { hello: 1, world: 1234 };
  const postHeaders = createSignature("POST", "/api/v1/address/deposit-address", postBody);
  await fetch(postUrl, {
    method: "POST",
    headers: postHeaders,
    body: JSON.stringify(postBody),
  })
    .then((response) => response.json())
    .then((data) => {
      console.log("API Response:", data);
    })
    .catch((error) => {
      console.error("Error:", error);
    });

  const getUrl = "http://127.0.0.1:3000/api/v1/address/deposit-address?depositLabel=poolETH1a";
  const getHeaders = createSignature("GET", "/api/v1/address/deposit-address");
  await fetch(getUrl, { method: "GET", headers: getHeaders })
    .then((response) => response.json())
    .then((data) => {
      console.log("API Response:", data);
    })
    .catch((error) => {
      console.error("Error:", error);
    });
};

test();
