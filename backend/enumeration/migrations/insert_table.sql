INSERT INTO properties (property_no, ownership_type, property_type, complex_name) VALUES
('PROP001', 'PRIVATE', 'RESIDENTIAL', 'Green Valley Apartments'),
('PROP002', 'PRIVATE', 'NON_RESIDENTIAL', NULL),
('PROP003', 'VACANT_LAND', 'MIXED', NULL),
('PROP004', 'CENTRAL_GOVERNMENT_50', 'RESIDENTIAL', 'Government Housing Complex'),
('PROP005', 'STATE_GOVERNMENT', 'NON_RESIDENTIAL', 'State Office Complex');


INSERT INTO construction_details (
    property_id, 
    floor_type, 
    wall_type, 
    roof_type, 
    wood_type
) VALUES (
    'c4d3926c-1e9a-4b3d-951e-2e919637dcfa',
    'RCC Slab',
    'Brick Wall',
    'Concrete Roof',
    'Teak Wood'
);

INSERT INTO gis_data (
    property_id, 
    source, 
    type, 
    entity_type
) VALUES (
    'c4d3926c-1e9a-4b3d-951e-2e919637dcfa', 
    'GPS',
    'POLYGON',
    'Residential Building'
);