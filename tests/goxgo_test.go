package goxgo

import (
	gxg "github.com/theplant/goxgo"
	"testing"
	"time"
)

func showResults(t *testing.T, results <-chan []string, done chan<- bool) {
	for {
		select {
		case res := <-results:
			t.Logf("result: %+v", res)
		}
	}
	done <- true
}

func TestGoXGo(t *testing.T) {
	dsn := gxg.DSN{
		Protocol: "tcp",
		Host:     "localhost",
		Port:     4242,
	}

	tokenizePayload := gxg.TokenizeRequest{
		Target: &gxg.CallTarget{Services: []string{"NLTK/tokenize"}, Version: "0.1"},
		Body:   "Give me a tokenized version of this unoptimzed body of text pls. Once successfully done we will try to stem the words too. Testing trying embodiment embodied",
		Locale: "en",
	}

	var tokenizeResponse gxg.TokenizeResponse
	gxg.Call(&dsn, &tokenizePayload, &tokenizeResponse)
	t.Logf("%+v", tokenizeResponse)

	stemPayload := gxg.StemRequest{
		Target: &gxg.CallTarget{Services: []string{"NLTK/stem"}, Version: "0.1"},
		Words:  tokenizeResponse.Tokens,
		Locale: tokenizeResponse.Locale,
	}

	var stemResponse gxg.StemResponse
	gxg.Call(&dsn, &stemPayload, &stemResponse)
	t.Logf("%+v", stemResponse)

	// goroutine requests
	var iter int
	iter = 3
	results := make(chan []string, iter)
	done := make(chan bool)
	for i := 0; i < iter; i++ {
		payload := gxg.TokenizeRequest{
			Target: &gxg.CallTarget{Services: []string{"NLTK/tokenize"}, Version: "0.1"},
			Body:   "Give me a tokenized version of this unoptimzed body of text pls. Once successfully done we will try to stem the words too. Testing trying embodiment embodied",
			Locale: "en",
		}
		go func() {
			var resp gxg.TokenizeResponse
			gxg.Call(&dsn, &payload, &resp)
			results <- resp.Tokens
		}()
	}
	go showResults(t, results, done)
	time.Sleep(0.5e9)

	// test vsm diff
	vsmComparePayload := gxg.VsmDiffRequest{
		Target:         &gxg.CallTarget{Services: []string{"NLTK/vsm_compare"}, Version: "0.1"},
		Drop_Stopwords: true,
		Stem_Words:     true,
		Docs: []string{`The domestic cat[1][2] (Felis catus[2] or Felis silvestris catus[4]) is a small, usually furry, domesticated, and carnivorous mammal. It is often called the housecat when kept as an indoor pet,[6] or simply the cat when there is no need to distinguish it from other felids and felines. Cats are valued by humans for companionship and their ability to hunt vermin and household pests.

Cats are similar in anatomy to the other felids, with strong, flexible bodies, quick reflexes, sharp retractable claws, and teeth adapted to killing small prey. Cat senses fit a crepuscular and predatory ecological niche. Cats can hear sounds too faint or too high in frequency for human ears, such as those made by mice and other small game. They can see in near darkness. Like most other mammals, cats have poorer color vision and a better sense of smell than humans.

Despite being solitary hunters, cats are a social species, and cat communication includes the use of a variety of vocalizations (meowing, purring, trilling, hissing, growling and grunting) as well as cat pheromones and types of cat-specific body language.[7]

Cats have a rapid breeding rate. Under controlled breeding, they can be bred and shown as registered pedigree pets, a hobby known as cat fancy. Failure to control the breeding of pet cats by spaying and neutering, and the abandonment of former household pets, has resulted in large numbers of feral cats worldwide, with a population of up to 60 million of these animals in the United States alone, requiring population control.[8]

Since cats were cult animals in ancient Egypt, they were commonly believed to have been domesticated there,[9] but there may have been instances of domestication as early as the Neolithic.[10]

A genetic study in 2007 revealed that domestic cats are descended from African wildcats (Felis silvestris lybica) c. 8000 BCE, in the Middle East.[9][11] According to Scientific American, cats are the most popular pet in the world, and are now found almost every place where people live.[12]`,
			`The domestic dog (Canis lupus familiaris)[2][3] is a subspecies of the gray wolf (Canis lupus), a member of the Canidae family of the mammalian order Carnivora. The term "domestic dog" is generally used for both domesticated and feral varieties. The dog may have been the first animal to be domesticated, and has been the most widely kept working, hunting, and pet animal in human history. The word "dog" may also mean the male of a canine species,[4] as opposed to the word "bitch" for the female of the species.

The present lineage of dogs was domesticated from gray wolves probably about 15,000 years ago.[5] Though remains of domesticated dogs have been found in Siberia and Belgium from about 33,000 years ago, none of those lineages seem to have survived the Last Glacial Maximum. Although DNA testing suggests an evolutionary split between dogs and wolves around 100,000 years ago, no fossil specimens prior to 33,000 years ago are clearly morphologically domesticated dog.[6][7][8]

Dogs' value to early human hunter-gatherers led to them quickly becoming ubiquitous across world cultures. Dogs perform many roles for people, such as hunting, herding, pulling loads, protection, assisting police and military, companionship, and, more recently, aiding handicapped individuals. This impact on human society has given them the nickname "Man's Best Friend" in the Western world. In some cultures, dogs are also a source of meat.[9][10] In 2001, there were estimated to be 400 million dogs in the world.[11]

Most breeds of dogs are at most a few hundred years old, having been artificially selected for particular morphologies and behaviors by people for specific functional roles. Through this selective breeding, the dog has developed into hundreds of varied breeds, and shows more behavioral and morphological variation than any other land mammal.[12] For example, height measured to the withers ranges from 15.2 centimetres (6.0 in) in the Chihuahua to about 76 cm (30 in) in the Irish Wolfhound; color varies from white through grays (usually called "blue") to black, and browns from light (tan) to dark ("red" or "chocolate") in a wide variation of patterns; coats can be short or long, coarse-haired to wool-like, straight, curly, or smooth.[13] It is common for most breeds to shed this coat.`},
	}

	// And now about something completely different but it also containes the words body and optimze

	var vsmDiffResponse gxg.VsmDiffResponse
	gxg.Call(&dsn, &vsmComparePayload, &vsmDiffResponse)
	t.Logf("Compare cats and dogs: \n%+v\n", vsmDiffResponse)

	vsmComparePayload = gxg.VsmDiffRequest{
		Target:         &gxg.CallTarget{Services: []string{"NLTK/vsm_compare"}, Version: "0.1"},
		Drop_Stopwords: true,
		Stem_Words:     true,
		Docs: []string{`The domestic cat[1][2] (Felis catus[2] or Felis silvestris catus[4]) is a small, usually furry, domesticated, and carnivorous mammal. It is often called the housecat when kept as an indoor pet,[6] or simply the cat when there is no need to distinguish it from other felids and felines. Cats are valued by humans for companionship and their ability to hunt vermin and household pests.

Cats are similar in anatomy to the other felids, with strong, flexible bodies, quick reflexes, sharp retractable claws, and teeth adapted to killing small prey. Cat senses fit a crepuscular and predatory ecological niche. Cats can hear sounds too faint or too high in frequency for human ears, such as those made by mice and other small game. They can see in near darkness. Like most other mammals, cats have poorer color vision and a better sense of smell than humans.

Despite being solitary hunters, cats are a social species, and cat communication includes the use of a variety of vocalizations (meowing, purring, trilling, hissing, growling and grunting) as well as cat pheromones and types of cat-specific body language.[7]

Cats have a rapid breeding rate. Under controlled breeding, they can be bred and shown as registered pedigree pets, a hobby known as cat fancy. Failure to control the breeding of pet cats by spaying and neutering, and the abandonment of former household pets, has resulted in large numbers of feral cats worldwide, with a population of up to 60 million of these animals in the United States alone, requiring population control.[8]

Since cats were cult animals in ancient Egypt, they were commonly believed to have been domesticated there,[9] but there may have been instances of domestication as early as the Neolithic.[10]

A genetic study in 2007 revealed that domestic cats are descended from African wildcats (Felis silvestris lybica) c. 8000 BCE, in the Middle East.[9][11] According to Scientific American, cats are the most popular pet in the world, and are now found almost every place where people live.[12]`,
			`Originally, an engine was a mechanical device that converted force into motion. Military devices such as catapults, trebuchets and battering rams are referred to as siege engines. The term "gin" as in cotton gin is recognised as a short form of the Old French word engin, in turn from the Latin ingenium, related to ingenious. Most devices in the industrial revolution were called engines, and this is where the steam engine gained its name.[citation needed]

The term motor was originally used to distinguish the new internal combustion engine-powered vehicles from earlier vehicles powered by steam engines, such as the steam roller and motor roller, but may be used to refer to any engine.[citation needed]

In modern usage, the term engine describes devices capable of performing mechanical work, as in the original steam engine. In most cases, the work is produced by exerting a torque or linear force, which operates other machinery that generates electricity, pumps water, or compresses gas. In the context of propulsion systems, an air-breathing engine is one that uses atmospheric air to oxidise the fuel rather than supplying an independent oxidizer, as in a rocket.

While colloquially interchanged, there is a difference between a motor and an engine. In common usage, an engine burns or otherwise consumes fuel, and is differentiated from an electric machine (i.e., electric motor) that derives power without changing the composition of matter.[3] A heat engine may also serve as a prime mover, a component that transforms the flow or changes in pressure of a fluid into mechanical energy.[4] An automobile powered by an internal combustion engine may make use of various motors and pumps, but ultimately all such devices derive their power from the engine. Another way of looking at it is that a motor receives power from an external source, and then converts it into mechanical energy, while an engine creates power from pressure (derived directly from the explosive force of combustion or other chemical reaction, or secondarily from the action of some such force on other substances such as air, water, or steam).[5]

Devices converting heat energy into motion are commonly referred to simply as engines.[6]`},
	}

	gxg.Call(&dsn, &vsmComparePayload, &vsmDiffResponse)
	t.Logf("Compare cats and engines: \n%+v\n", vsmDiffResponse)
}
