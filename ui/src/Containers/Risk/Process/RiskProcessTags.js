import React, { useState } from 'react';

import Tags from 'Components/Tags';

const defaultTags = ['spicy', 'mild', 'bland'];

const RiskProcessTags = () => {
    const [tags, setTags] = useState(defaultTags);

    return <Tags type="Process" tags={tags} onChange={setTags} defaultOpen />;
};

export default RiskProcessTags;
