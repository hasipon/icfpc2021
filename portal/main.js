const dotenv = require('dotenv');
const puppeteer = require('puppeteer');
const fs = require('fs');

async function main() {
    console.log(process.cwd());
    // load .env file into process.env
    dotenv.config();
    const email = process.env.ICFPC_EMAIL;
    const password = process.env.ICFPC_PASSWORD;

    puppeteer.launch({
        headless: true,
	args: ['--no-sandbox', '--disable-setuid-sandbox'],
    }).then(async browser => {
        const page = await browser.newPage();

        if (fs.existsSync("cookies.json")) {
            const cookies = JSON.parse(fs.readFileSync("cookies.json", 'utf-8'));
            for (let cookie of cookies) {
                await page.setCookie(cookie);
            }
        } else {
            await login(page, email, password);
            const cookies = await page.cookies();
            fs.writeFileSync("cookies.json", JSON.stringify(cookies));
        }
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
