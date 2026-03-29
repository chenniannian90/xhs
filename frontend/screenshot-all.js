import puppeteer from 'puppeteer';
import fs from 'fs';

const pages = [
  { url: 'http://localhost:5173/register', name: '注册页面' },
  { url: 'http://localhost:5173/dashboard', name: '仪表板' },
  { url: 'http://localhost:5173/dashboard/categories', name: '分类管理' },
  { url: 'http://localhost:5173/dashboard/search', name: '搜索页面' },
  { url: 'http://localhost:5173/dashboard/settings', name: '设置页面' },
];

async function takeScreenshots() {
  const browser = await puppeteer.launch({
    headless: true,
    executablePath: '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome',
  });

  for (const page of pages) {
    console.log(`Screenshotting: ${page.name}...`);
    const p = await browser.newPage();
    await p.setViewport({ width: 1920, height: 1080 });
    await p.goto(page.url, { waitUntil: 'networkidle0' });
    await new Promise(resolve => setTimeout(resolve, 2000));

    const filename = `screenshot-${page.name.replace(/\s+/g, '-')}.png`;
    await p.screenshot({
      path: filename,
      fullPage: true,
    });

    fs.copyFile(filename, `/Users/mac-new/work/navhub/${filename}`, (err) => {
      if (err) console.error('Error copying:', err);
    });

    console.log(`  ✓ Saved to ${filename}`);
    await p.close();
  }

  await browser.close();
  console.log('\nAll screenshots complete!');
}

takeScreenshots().catch(console.error);
