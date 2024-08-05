import { useState } from 'react';

// import {
//   // DataFrame,
//   // DataLink,
//   // VariableSuggestion
// } from '@grafana/data';

import { Button } from '../Button';
import { Modal } from '../Modal/Modal';

import { Calculation, CalculationEditor } from './CalculationEditor';
// import { DataLinkEditor } from 'DataLinkEditor';

interface CalculationEditorModalContentProps {
  link: Calculation;
  // index: number;
  // data: DataFrame[];
  // getSuggestions: () => VariableSuggestion[];
  // onSave: (index: number, ink: DataLink) => void;
  // onCancel: (index: number) => void;
}

export const CalculationEditorModalContent = ({
  link,
  // index,
  // getSuggestions,
  // onSave,
  // onCancel,
}: CalculationEditorModalContentProps) => {
  const [
    dirtyLink,
    // setDirtyLink
  ] = useState(link);
  return (
    <>
      <CalculationEditor
        value={dirtyLink}
        // index={index}
        // isLast={false}
        // suggestions={getSuggestions()}
        // onChange={(
        //   // index,
        //   link
        // ) => {
        //   setDirtyLink(link);
        // }}
      />
      <Modal.ButtonRow>
        <Button
          variant="secondary"
          // onClick={() => onCancel(index)}
          fill="outline"
        >
          Cancel
        </Button>
        <Button
        // onClick={() => {
        //   onSave(index, dirtyLink);
        // }}
        >
          Save
        </Button>
      </Modal.ButtonRow>
    </>
  );
};
