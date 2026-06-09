require("dotenv").config();
const axios = require("axios");

const { App } = require("@slack/bolt");

const app = new App({
    token: process.env.SLACK_BOT_TOKEN,
    appToken: process.env.SLACK_APP_TOKEN,
    socketMode: true
})

app.command("/d4lm-ping", async ({ command, ack, respond}) => {
    const start = Date.now()
    await ack();
    const latency = Date.now - start;
    await respond({ text: `Pong!\nLatency: ${latency}ms` });
});

app.command("/d4lm-help", async ({ command, ack, respond}) => {
   await ack();
   await respond({
       text:
           'Available Commands: \n' +
           '/d4lm-ping - Check bot latency'
   });
});

app.command("/d4lm-catfact", async({ ack, respond }) => {
    await ack();

    try {
        const response = await axios.get("https://catfact.ninja/fact");
        await respond({ text: `Cat Fact: ${response.data.fact}` });
    } catch (err) {
        await respond({ text: "Failed to fetch a cat fact." });
    }
});

app.command("/d4lm-joke", async ({ ack, respond}) => {
   await ack();

   try {
       const response = await axios.get("https://official-joke-api.appspot.com/random_joke");
       await respond({
           text: `${response.data.setup} \n${response.data.punchline}`
       });
   } catch (err) {
       await respond({ text: "Failed to fetch a joke." })
   }
});

(async () => {
    await app.start();
    console.log("bot is running!");
})();