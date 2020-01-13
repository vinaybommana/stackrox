import React, { useState } from 'react';

import Tags from 'Components/Tags';

const defaultTags = ['spicy', 'mild', 'bland'];

const ViolationTags = () => {
    const [tags, setTags] = useState(defaultTags);

    return <Tags type="Violation" tags={tags} onChange={setTags} defaultOpen />;
};

export default React.memo(ViolationTags);
