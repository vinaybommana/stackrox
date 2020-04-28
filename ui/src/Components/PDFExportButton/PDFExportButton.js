import React from 'react';
import PropTypes from 'prop-types';
import html2canvas from 'html2canvas';
import { FileText } from 'react-feather';
import logError from 'modules/logError';
import { toast } from 'react-toastify';

import { getDate, addBrandedTimestampToString } from 'utils/dateUtils';
import Button from 'Components/Button';
import {
    createPDFContainerElement,
    createPDFHeaderElement,
    createPDFBodyElement,
    computeStyles,
    addElementToRootNode,
    removeElementFromRootNode,
    savePDF,
} from './pdfUtils';

const PDFExportButton = ({ fileName, pdfId, startExportingPDF, finishExportingPDF }) => {
    function exportPDF() {
        // This hides all the pdf generation behind an exporting screen
        startExportingPDF();

        const pdfTitle = `StackRox ${fileName}`;
        const currentTimestamp = getDate(new Date());
        const pdfFileName = addBrandedTimestampToString(fileName);

        // creates a container element that will include everything necessary to convert to a PDF
        const pdfContainerElement = createPDFContainerElement();

        // add the StackRox header to the container element
        const pdfHeaderElement = createPDFHeaderElement(pdfTitle, currentTimestamp);
        pdfContainerElement.appendChild(pdfHeaderElement);

        // create a clone of the element to be exported and add it to the body of the container
        const pdfBodyElement = createPDFBodyElement();
        const elementToBeExported = document.getElementById(pdfId);
        const clonedElementToBeExported = elementToBeExported.cloneNode(true);
        pdfBodyElement.appendChild(clonedElementToBeExported);
        pdfContainerElement.appendChild(pdfBodyElement);

        // we need to add the container element to the DOM in order to compute the styles and eventually convert it from HTML -> Canvas -> PNG -> PDF
        addElementToRootNode(pdfContainerElement);

        // we need to compute styles into inline styles in order for html2canvas to properly work
        computeStyles(pdfBodyElement);

        // convert HTML -> Canvas
        html2canvas(pdfContainerElement, {
            scale: 1,
            allowTaint: true,
        })
            .then((canvas) => {
                // convert Canvas -> PNG -> PDF
                savePDF(canvas, pdfFileName);
                // Remember to clean up after yourself. This makes sure to remove any added elements to the DOM after they're used
                removeElementFromRootNode(pdfContainerElement);
                // remove the exporting screen
                finishExportingPDF();
            })
            .catch((error) => {
                logError(error);
                finishExportingPDF();
                toast('An error occurred while exporting. Please try again.');
            });
    }

    return (
        <Button
            className="btn btn-base"
            icon={<FileText className="h-4 w-4 mx-2" />}
            text="Export"
            onClick={exportPDF}
        />
    );
};

PDFExportButton.propTypes = {
    fileName: PropTypes.string.isRequired,
    pdfId: PropTypes.string.isRequired,
    startExportingPDF: PropTypes.func.isRequired,
    finishExportingPDF: PropTypes.func.isRequired,
};

export default PDFExportButton;
