const dotenv = require('dotenv');
const puppeteer = require('puppeteer');
const fs = require('fs');

function main() {
    console.log(process.cwd());
    // load .env file into process.env
    dotenv.config();
    const email = process.env.ICFPC_EMAIL;
    const password = process.env.ICFPC_PASSWORD;
    console.log(email);

    puppeteer.launch({
        headless: true,
    }).then(async browser => {
        const page = await browser.newPage();
        await login(page, email, password);
        await save_problems_json(page);
        await browser.close();
    });
}

async function login(page, email, password) {
    await Promise.all([
        page.goto('https://poses.live/login'),
        page.waitForNavigation({waitUntil: 'domcontentloaded'})
    ]);
    await page.type('#login\\.email', email);
    await page.type('#login\\.password', password);
    await Promise.all([
        page.click('body > section > form > input[type=submit]:nth-child(6)'),
        page.waitForNavigation({waitUntil: 'domcontentloaded'})
    ]);
}

async function save_problems_json(page) {
    await Promise.all([
        page.goto('https://poses.live/problems'),
        page.waitForNavigation({waitUntil: 'domcontentloaded'})
    ]);

    const result = await page.evaluate(() => {
        const rows = document.querySelectorAll('body > section > table > tbody > tr');
        return Array.from(rows, row => {
            const columns = row.querySelectorAll('td');
            return Array.from(columns, column => column.innerText);
        });
    });

    fs.writeFileSync("problems.json", JSON.stringify(result.filter((e) => 0 < e.length)));
}

main();