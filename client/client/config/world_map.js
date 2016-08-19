WorldMap = [
    {
        x: 3,
        y: -5,
        is_explored: true,
        "terrains": {
            "forest": {
                "type": "",
                "quality": 4,
                "passability": 0.5,
                "percentArea": 100.0
            }
        }
    },
    {
        x: 3,
        y: -6,
        is_explored: true,
        "terrains": {
            "forest": {
                "type": "",
                "quality": 3,
                "passability": 0.7,
                "percentArea": 95.0
            },
            "meadow": {
                "type": "",
                "quality": 1,
                "passability": 0.9,
                "percentArea": 5.0
            }
        }
    },
    {
        x: 3,
        y: -7,
        is_explored: true,
        "terrains": {
            "roads": [
                {
                    "size": 1,
                    "onRiver": {
                        "0": "bottom"
                    },
                    "direction": "W-E"
                }
            ],
            "forest": {
                "type": "",
                "quality": 2,
                "passability": 0.7,
                "percentArea": 70.0
            },
            "meadow": {
                "type": "",
                "quality": 2,
                "passability": 0.8,
                "percentArea": 29.9
            },
            "rivers": [
                {
                    "size": 1,
                    "bridge": 0,
                    "onRoad": {
                        "0": "top"
                    },
                    "quality": 1,
                    "direction": "W-E",
                    "percentArea": 0.1
                }
            ]
        }
    },
    {
        x: 3,
        y: -8,
        is_explored: true,
        "terrains": {
            "field": {
                "type": "",
                "quality": 3,
                "passability": 0.8,
                "percentArea": 90.0
            },
            "forest": {
                "type": "",
                "quality": 1,
                "passability": 0.9,
                "percentArea": 5
            },
            "meadow": {
                "type": "",
                "quality": 1,
                "passability": 0.9,
                "percentArea": 5.0
            }
        }
    },
    {
        x: 3,
        y: -9,
        is_explored: true,
        "terrains": {
            "hill": {
                "type": "",
                "quality": 2,
                "passability": 0.7,
                "percentArea": 80.0
            },
            "forest": {
                "type": "",
                "quality": 1,
                "passability": 0.9,
                "percentArea": 10
            },
            "meadow": {
                "type": "",
                "quality": 1,
                "passability": 0.9,
                "percentArea": 10.0
            }
        }
    },
    {
        x: 4,
        y: -5,
        is_explored: true,
        "terrains": {"forest": {"type": "", "quality": 4, "passability": 0.5, "percentArea": 100.0}}
    },
    {
        x: 4,
        y: -6,
        is_explored: true,
        "terrains": {
            "forest": {"type": "", "quality": 3, "passability": 0.7, "percentArea": 95.0},
            "meadow": {"type": "", "quality": 1, "passability": 0.9, "percentArea": 5.0}
        }
    },
    {
        x: 4,
        y: -7,
        is_explored: true,
        "terrains": {
            "roads": [{"size": 1, "onRiver": {"0": "bottom"}, "direction": "W-E"}],
            "forest": {"type": "", "quality": 2, "passability": 0.7, "percentArea": 70.0},
            "meadow": {"type": "", "quality": 2, "passability": 0.8, "percentArea": 29.9},
            "rivers": [{
                "size": 1,
                "bridge": 0,
                "onRoad": {"0": "top"},
                "quality": 1,
                "direction": "W-E",
                "percentArea": 0.1
            }]
        }
    },
    {
        x: 4,
        y: -8,
        is_explored: true,
        "terrains": {
            "field": {"type": "", "quality": 3, "passability": 0.8, "percentArea": 90.0},
            "forest": {"type": "", "quality": 1, "passability": 0.9, "percentArea": 5},
            "meadow": {"type": "", "quality": 1, "passability": 0.9, "percentArea": 5.0}
        }
    },
    {
        x: 4,
        y: -9,
        is_explored: true,
        "terrains": {
            "hill": {"type": "", "quality": 2, "passability": 0.7, "percentArea": 80.0},
            "forest": {"type": "", "quality": 1, "passability": 0.9, "percentArea": 10},
            "meadow": {"type": "", "quality": 1, "passability": 0.9, "percentArea": 10.0}
        }
    },
    {
        x: 5,
        y: -5,
        is_explored: true,
        "terrains": {
            "roads": [{"size": 2, "direction": "S-N"}],
            "forest": {"type": "", "quality": 4, "passability": 0.5, "percentArea": 100.0}
        }
    },
    {
        x: 5,
        y: -6,
        is_explored: true,
        "terrains": {
            "roads": [{"size": 2, "direction": "S-N"}],
            "forest": {"type": "", "quality": 3, "passability": 0.7, "percentArea": 95.0},
            "meadow": {"type": "", "quality": 1, "passability": 0.9, "percentArea": 5.0}
        }
    },
    {
        x: 5,
        y: -7,
        is_explored: true,
        "terrains": {
            "roads": [{"size": 1, "onRiver": {"0": "bottom"}, "direction": "W-C"}, {
                "size": 2,
                "onRiver": {"0": "across"},
                "direction": "S-N"
            }],
            "urban": {
                "type": "village",
                "onRoad": {"0": "across", "1": "across"},
                "townId": "12314",
                "onRiver": {"0": "across"},
                "percentArea": 15
            },
            "forest": {"type": "", "quality": 1, "passability": 0.8, "percentArea": 15},
            "meadow": {"type": "", "quality": 1, "passability": 0.9, "percentArea": 69.9},
            "rivers": [{
                "size": 1,
                "bridge": 3,
                "onRoad": {"0": "top", "1": "across"},
                "quality": 1,
                "direction": "W-E",
                "percentArea": 0.1
            }]
        }
    },
    {
        x: 5,
        y: -8,
        is_explored: true,
        "terrains": {
            "field": {"type": "", "quality": 2, "passability": 0.9, "percentArea": 45.0},
            "roads": [{"size": 2, "direction": "S-N"}],
            "forest": {"type": "", "quality": 1, "passability": 0.9, "percentArea": 5},
            "meadow": {"type": "", "quality": 1, "passability": 0.9, "percentArea": 50.0}
        }
    },
    {
        x: 5,
        y: -9,
        is_explored: true,
        "terrains": {
            "hill": {"type": "", "quality": 2, "passability": 0.7, "percentArea": 60.0},
            "roads": [{"size": 2, "direction": "S-N"}],
            "forest": {"type": "", "quality": 1, "passability": 0.9, "percentArea": 10},
            "meadow": {"type": "", "quality": 2, "passability": 0.9, "percentArea": 30.0}
        }
    },
    {
        x: 6,
        y: -5,
        is_explored: true,
        "terrains": {
            "swamp": {"type": "", "quality": 2, "passability": 0.4, "percentArea": 10.0},
            "forest": {"type": "", "quality": 4, "passability": 0.5, "percentArea": 85.0},
            "rivers": [{"size": 5, "bridge": 0, "quality": 3, "direction": "E-N", "percentArea": 5.0}]
        }
    },
    {
        x: 6,
        y: -6,
        is_explored: true,
        "terrains": {
            "forest": {"type": "", "quality": 3, "passability": 0.7, "percentArea": 95.0},
            "meadow": {"type": "", "quality": 2, "passability": 0.9, "percentArea": 5.0}
        }
    },
    {
        x: 6,
        y: -7,
        is_explored: true,
        "terrains": {
            "forest": {"type": "", "quality": 1, "passability": 0.8, "percentArea": 20},
            "meadow": {"type": "", "quality": 2, "passability": 0.9, "percentArea": 79.9},
            "rivers": [{"size": 1, "bridge": 0, "quality": 1, "direction": "W-E", "percentArea": 0.1}]
        }
    },
    {
        x: 6,
        y: -8,
        is_explored: true,
        "terrains": {
            "forest": {"type": "", "quality": 1, "passability": 0.9, "percentArea": 5},
            "meadow": {"type": "", "quality": 3, "passability": 0.8, "percentArea": 95.0}
        }
    },
    {
        x: 6,
        y: -9,
        is_explored: true,
        "terrains": {
            "forest": {"type": "", "quality": 1, "passability": 0.9, "percentArea": 5},
            "meadow": {"type": "", "quality": 3, "passability": 0.8, "percentArea": 95.0}
        }
    },
    {
        x: 7,
        y: -5,
        is_explored: true,
        "terrains": {
            "swamp": {"type": "", "quality": 3, "passability": 0.3, "percentArea": 30.0},
            "forest": {"type": "", "quality": 4, "passability": 0.5, "percentArea": 65.0},
            "rivers": [{"size": 5, "bridge": 0, "quality": 3, "direction": "S-W", "percentArea": 5.0}]
        }
    },
    {
        x: 7,
        y: -6,
        is_explored: true,
        "terrains": {
            "swamp": {"type": "", "quality": 1, "passability": 0.4, "percentArea": 5.0},
            "forest": {"type": "", "quality": 4, "passability": 0.5, "percentArea": 90.0},
            "rivers": [{"size": 5, "bridge": 0, "quality": 3, "direction": "S-N", "percentArea": 5.0}]
        }
    },
    {
        x: 7,
        y: -7,
        is_explored: true,
        "terrains": {
            "forest": {"type": "", "quality": 3, "passability": 0.6, "percentArea": 90.0},
            "meadow": {"type": "", "quality": 2, "passability": 0.9, "percentArea": 4.0},
            "rivers": [{"size": 5, "bridge": 0, "quality": 3, "direction": "S-N", "percentArea": 5.0}, {
                "size": 1,
                "bridge": 0,
                "quality": 1,
                "direction": "W-C",
                "percentArea": 1.0
            }]
        }
    },
    {
        x: 7,
        y: -8,
        is_explored: true,
        "terrains": {
            "forest": {"type": "", "quality": 3, "passability": 0.6, "percentArea": 90.0},
            "meadow": {"type": "", "quality": 2, "passability": 0.9, "percentArea": 5.0},
            "rivers": [{"size": 5, "bridge": 0, "quality": 3, "direction": "S-N", "percentArea": 5.0}]
        }
    },
    {
        x: 7,
        y: -9,
        is_explored: true,
        "terrains": {
            "forest": {"type": "", "quality": 3, "passability": 0.6, "percentArea": 90.0},
            "meadow": {"type": "", "quality": 2, "passability": 0.9, "percentArea": 5.0},
            "rivers": [{"size": 5, "bridge": 0, "quality": 3, "direction": "S-N", "percentArea": 5.0}]
        }
    }
];