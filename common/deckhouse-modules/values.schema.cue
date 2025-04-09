
#WeekDay: "Mon" | "Tue" | "Wed" | "Thu" | "Fri" | "Sat" | "Sun"

#Window: {
		to: string
		from: string
		days: [...#WeekDay]
}

#DeckhouseVal: {
		updateMode?: "Auto" | "Manual"
		releaseChannel?: "Stable" | "Beta" | "Alpha"
		version?: number
		windows: [...#Window]
}

clusterName: string
clusterDomain: string

deckhouse!: #DeckhouseVal
