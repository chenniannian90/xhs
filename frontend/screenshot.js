import puppeteer from 'puppeteer';
import fs from 'fs';

async function takeScreenshot(url = 'http://localhost:5173/', outputFile = 'screenshot.png') {
  console.log('Launching browser...');
  const browser = await puppeteer.launch({
    headless: true,
    executablePath: '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome',
  });

  console.log(`Navigating to ${url}...`);
  const page = await browser.newPage();

  // Set viewport size
  await page.setViewport({ width: 1920, height: 1080 });

  // Navigate to URL
  await page.goto(url, { waitUntil: 'networkidle0' });

  // Wait a bit for any animations
  await new Promise(resolve => setTimeout(resolve, 3000));

  console.log('Taking screenshot...');
  await page.screenshot({
    path: outputFile,
    fullPage: true,
  });

  await browser.close();
  console.log(`Screenshot saved to ${outputFile}`);

  // Also save to project root for easy access
  const projectPath = '/Users/mac-new/work/navhub/screenshot.png';
  fs.copyFile(outputFile, projectPath, (err) => {
    if (err) console.error('Error copying file:', err);
    else console.log(`Also saved to ${projectPath}`);
  });
}

takeScreenshot().catch(console.error);
