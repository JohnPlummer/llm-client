package scorer

const maxBatchSize = 10

const batchScorePrompt = `Score each of the following Reddit post titles. Consider these categories:
- Regular venues (restaurants, bars, cafes, museums, galleries, etc.)
- Local attractions and points of interest
- Entertainment events (music, theatre, comedy, sports, etc.)
- Cultural events and festivals
- Markets and shopping areas
- Parks and outdoor spaces
- Family-friendly activities
- Seasonal or special events
- Hidden gems and local recommendations

Scoring guidelines:
90-100: Title directly references specific venues, events, or activities
70-89: Title suggests discussion of activities or places
40-69: Title might contain some relevant information
1-39: Title has low probability of relevant information
0: Title clearly indicates no relevant activity information

Posts to score:
%s`
