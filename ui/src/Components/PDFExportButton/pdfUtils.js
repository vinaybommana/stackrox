import StackroxLogo from 'images/stackrox-logo.png';
import computedStyleToInlineStyle from 'computed-style-to-inline-style';
import JSPDF from 'jspdf';

/**
 * Creates a container div HTML element that will wrap around all the content to be exported
 * @returns {HTMLElement}
 */
export function createPDFContainerElement() {
    const pdfContainer = document.createElement('div');
    pdfContainer.id = 'pdf-container';
    pdfContainer.className = 'flex flex-1 flex-col h-full -z-1 absolute top-0 left-0 theme-light';
    return pdfContainer;
}

/**
 * Creates a header HTML element that will contain the StackRox logo, PDF title, and the current time
 *  @param {string} pdfTitle - The title to display in the top right section of the header
 *  @param {string} timestamp - The timestamp to display in the top right section of the header
 *  @returns {HTMLElement}
 */
export function createPDFHeaderElement(pdfTitle, timestamp) {
    const div = `<div class="theme-light flex justify-between bg-primary-800 items-center text-primary-100 h-32">
            <img alt="stackrox-logo" src=${StackroxLogo} class="h-24" />
            <div class="pr-4 text-right">
                <div class="text-2xl">${pdfTitle}</div>
                <div class="pt-2 text-xl">${timestamp}</div>
            </div>
        </div>`;
    const header = document.createElement('header');
    header.id = 'pdf-header';
    header.innerHTML = div;
    return header;
}

/**
 * Creates a div HTML element that will contain the content being exported
 * @returns {HTMLElement}
 */
export function createPDFBodyElement() {
    const body = document.createElement('div');
    body.id = 'pdf-body';
    body.className = 'flex flex-1 border-b border-base-300 -z-1';
    return body;
}

/**
 * Converts an HTML element's computed CSS to inline CSS
 * @param {HTMLElement} element
 */
export function computeStyles(element) {
    const isThemeDark = document.body.className.includes('theme-dark');

    // if dark mode is enabled, we want to switch to light mode for exporting to PDF
    if (isThemeDark) {
        document.body.classList.remove('theme-dark');
        document.body.classList.add('theme-light');
    }

    computedStyleToInlineStyle(element, {
        recursive: true,
        properties: ['width', 'height', 'fill', 'style', 'class', 'stroke', 'font', 'font-size'],
    });

    // if dark mode was previously enabled, we want to switch back after styles are computed
    if (isThemeDark) {
        document.body.classList.remove('theme-light');
        document.body.classList.add('theme-dark');
    }
}

/**
 * Adds an element to the Root Node
 *  @param {HTMLElement} element
 */
export function addElementToRootNode(element) {
    document.getElementById('root').appendChild(element);
}

/**
 * Removes an element from the Root Node
 *  @param {HTMLElement} element
 */
export function removeElementFromRootNode(element) {
    if (element?.parentNode) {
        element.parentNode.removeChild(element);
    }
}

/**
 *  Converts a Canvas element -> PNG -> PDF
 *  @param {HTMLElement} canvas
 *  @param {string} pdfFileName - The PDF file name
 */
export function savePDF(canvas, pdfFileName) {
    const pdf = new JSPDF();
    const imgData = canvas.toDataURL('image/png');

    // we want the width to be 100% of the PDF page, but the height to scale within the w/h ratio of the Canvas element
    const imgProps = pdf.getImageProperties(imgData);
    const pdfWidth = pdf.internal.pageSize.getWidth();
    const pdfHeight = (imgProps.height * pdfWidth) / imgProps.width;

    pdf.addImage(imgData, 'PNG', 0, 0, pdfWidth, pdfHeight);
    pdf.save(pdfFileName);
}
