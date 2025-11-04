import os
import uuid
from playwright.async_api import async_playwright

OUTPUT_DIR = "temp_pdfs"
os.makedirs(OUTPUT_DIR, exist_ok=True)

async def converter_html_para_pdf(html_content: str) -> str:
    pdf_filename = f"{uuid.uuid4()}.pdf"
    pdf_path = os.path.join(OUTPUT_DIR, pdf_filename)
    
    async with async_playwright() as p:
        browser = await p.chromium.launch()
        page = await browser.new_page()
        

        await page.set_content(html_content, wait_until="commit")
        
        await page.wait_for_timeout(2000)

        await page.pdf(
            path=pdf_path,
            format="A4",
            print_background=True,
            margin={"top": "20mm", "bottom": "20mm", "left": "20mm", "right": "20mm"}
        )
        
        await browser.close()
        
    return pdf_path