import React, { useState } from 'react';

import Tags from 'Components/Tags';

const defaultTags = ['spicy', 'mild', 'bland'];

const ProcessTags = () => {
    const [tags, setTags] = useState(defaultTags);

    return <Tags type="Process" tags={tags} onChange={setTags} defaultOpen />;
};

export default React.memo(ProcessTags);
