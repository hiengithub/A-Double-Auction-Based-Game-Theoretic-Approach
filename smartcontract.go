

package main


import (
    "bytes"
    "encoding/json"
    "fmt"
    "strconv"
    "strings"
    // "math/rand"
    // "math"
    "sort"
    "time"

    "github.com/hyperledger/fabric-chaincode-go/shim"
    sc "github.com/hyperledger/fabric-protos-go/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type UserInfo struct {
    XMax [13]float64        `json:"xMax"`
    SurplusEnergy [13]float64        `json:"surplusEnergy"`
    Prices [13]float64        `json:"prices"`
    Identify string     `json:"identify"`
    ClearQuantity [13]float64     `json:"clearQuantity"`
    ClearPrice [13]float64    `json:"clearPrice"`
}

type ParseUserInfo struct {
    XMax [13]string        `json:"xMax"`
    SurplusEnergy [13]string        `json:"surplusEnergy"`
    Prices [13]string        `json:"prices"`
    Identify string     `json:"identify"`
    ClearQuantity [13]string     `json:"clearQuantity"`
    ClearPrice [13]string    `json:"clearPrice"`
}

type trackedModel struct {
    PMax string        `json:"pMax"`
    PMin string        `json:"pMin"`
    Price string       `json:"price"`
    BuyWinners []string   `json:"buyWinners"`
    SellWinners []string   `json:"sellWinners"`
    OptimalPrice string   `json:"optimalPrice"`
    OptimalQuantity []float64   `json:"optimalQuantity"`
    MaximumSW float64   `json:"maximumSW"`
    OfferPrice []string   `json:"offerPrice"`
    TotalQuantitySeller float64     `json:"totalQuantitySeller"`
    TotalQuantityBuyer float64     `json:"totalQuantityBuyer"`
}

type ResponseModel struct {
    UserID string        `json:"userID"`
    ClearQuantity string        `json:"clearQuantity"`
    ClearPrice string        `json:"clearPrice"`
    Role string        `json:"role"`
}

const END_KEY = "EndKey"
const TRACKED_MODEL_KEY = "TrackedModel"

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
    return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

    // Retrieve the requested Smart Contract function and arguments
    function, args := APIstub.GetFunctionAndParameters()
    // Route to the appropriate handler function to interact with the ledger appropriately
    if function == "startGame" {
        return s.startGame(APIstub,args)
    }  
    if function == "emptyContract" {
        return s.emptyContract(APIstub)
    } 
    if function == "initLedger" {
        return s.initLedger(APIstub,args)
    }

    return shim.Error("Invalid Smart Contract function name.")
}


func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    // fmt.Println(args[0]) 
    // c := make(map[string]string)
    var userList []ParseUserInfo
    _ = json.Unmarshal([]byte(args[0]), &userList)
    fmt.Println("total number of buyers and seller: ",len(userList)) 
    for idx := range userList {
        var user UserInfo
        
        for idx, value := range userList[idx].XMax {
            xmax64, _ := strconv.ParseFloat(value, 64)
            user.XMax[idx] = xmax64
        }

        for idx, value := range userList[idx].SurplusEnergy {
            surplusEnergy64, _ := strconv.ParseFloat(value, 64)
            user.SurplusEnergy[idx] = surplusEnergy64
        }

        for idx, value := range userList[idx].Prices {
            price64, _ := strconv.ParseFloat(value, 64)
            user.Prices[idx] = price64
        }
    
        // user.Prices = prices64
        user.Identify = userList[idx].Identify
        user.ClearPrice = [13]float64{0,0,0,0,0,0,0,0,0,0,0,0,0}
        user.ClearQuantity = [13]float64{0,0,0,0,0,0,0,0,0,0,0,0,0}
        playerAsBytes, _ := json.Marshal(user)
        err := APIstub.PutState(userList[idx].Identify, playerAsBytes)
        if err != nil {
            fmt.Println("err to PutState ",user.Identify," is ", err)        
        }

        // fmt.Println("\nIdentify",user.Identify)
        // fmt.Println("XMax",user.XMax)
        // fmt.Println("SurplusEnergy",user.SurplusEnergy)
        // fmt.Println("Prices",user.Prices)
        // fmt.Println("ClearPrice",user.ClearPrice)
        // fmt.Println("ClearQuantity",user.ClearQuantity)
    }

    APIstub.DelState(TRACKED_MODEL_KEY)
    playerAsBytes, _ := json.Marshal(args[1])
    APIstub.PutState(END_KEY, playerAsBytes)

    return shim.Success(nil)
}

func findWinner(APIstub shim.ChaincodeStubInterface,t int) ([]UserInfo,[]UserInfo,int,int){ 

    var Buyers []UserInfo
    var Sellers []UserInfo

    var endKey string
    endKeyBytes,_ := APIstub.GetState(END_KEY)
    json.Unmarshal(endKeyBytes, &endKey)

    // fmt.Println("EndKey",endKey)
    resultsIterator, _ := APIstub.GetStateByRange("user0001", endKey)
    defer resultsIterator.Close()

    for resultsIterator.HasNext() {
        queryResponse, _ := resultsIterator.Next()
        user := UserInfo{}
        json.Unmarshal(queryResponse.Value, &user)
        // fmt.Println("\nIdentify:", user.Identify)
        if user.SurplusEnergy[t] > 0 {
            Sellers = append(Sellers,user)
        } else {
            Buyers = append(Buyers,user)
        }
    }
    // fmt.Println("total number of Buyers ",len(Buyers)," ,total number of Sellers ",len(Sellers))
    sort.SliceStable(Buyers, func(i, j int) bool { return Buyers[i].Prices[t] > Buyers[j].Prices[t] })
    sort.SliceStable(Sellers, func(i, j int) bool { return Sellers[i].Prices[t] < Sellers[j].Prices[t] })
    

    // fmt.Println(Buyers, Sellers) 
    if len(Buyers) > 0 && len(Sellers) > 0 {
        i,j := optimal_trade(Buyers,Sellers,t)
        return Buyers,Sellers,i,j
    } else {
        return Buyers,Sellers,-1,-1
    }
    
}

func optimal_trade(Buyers []UserInfo,Sellers []UserInfo,t int) (int,int) {
    j := 0 
    i := 0 
    // breakEvenIndex = (0,0)
    cumBuyersClearedQuantity := Buyers[i].XMax[t]
    cumSellersClearedQuantity := Sellers[j].SurplusEnergy[t]

    for i<len(Buyers) && j<len(Sellers) {
        buyer := Buyers[i]
        seller := Sellers[j]

        if cumBuyersClearedQuantity < cumSellersClearedQuantity {
            if i + 1 < len(Buyers) {
                i += 1
            } else {
                break
            }
            if Buyers[i].Prices[t] <= seller.Prices[t]{
                // breakEvenIndex = (i-1,j)
                return i-1,j
            }
            cumBuyersClearedQuantity += Buyers[i].XMax[t]
        } else if cumBuyersClearedQuantity > cumSellersClearedQuantity {
            if j + 1 < len(Sellers) {
                j += 1
            } else {
                break
            }
            if buyer.Prices[t] <= Sellers[j].Prices[t] {
                // breakEvenIndex = (i,j-1)
                return i,j-1
            }
            cumSellersClearedQuantity += Sellers[j].SurplusEnergy[t]
        } else {
            if j + 1 < len(Sellers) {
                j += 1
            }
            if i + 1 < len(Buyers) {
                i += 1
            }
            if Buyers[i].Prices[t] <= Sellers[j].Prices[t] {
                return i-1,j-1
            }
            cumSellersClearedQuantity += Sellers[j].SurplusEnergy[t]
            cumBuyersClearedQuantity += Buyers[i].XMax[t]
        }
        // breakEvenIndex = (i,j)
    }

    return i,j
}

func (s *SmartContract) startGame(APIstub shim.ChaincodeStubInterface,args []string) sc.Response {
    // _ := args[0]
    start := time.Now()
    t,_ := strconv.Atoi(args[2])
    increment := 0.1
    trackedValue,_ := APIstub.GetState(TRACKED_MODEL_KEY)
    // p,_ := APIstub.GetState("price")

    if trackedValue == nil { 
        Buyers,Sellers,i,j := findWinner(APIstub,t)

        // fmt.Println("before delete len(Buyers),len(Sellers) ",len(Buyers),len(Sellers)," ,i,j",i,j)
        var pmax float64
        if i > 0 {
            pmax = Buyers[i].Prices[t]
            Buyers = Buyers[:i]
            // del Buyers[i::]
            i -= 1
        } else {
            pmax = Buyers[0].Prices[t]
            if len(Buyers) > 1 {
                Buyers = Buyers[:i+1]    
            }
            // del Buyers[i+1::]
        }
        if j > 0 {
            Sellers = Sellers[:j]
            // del Sellers[j::]
            j -= 1
        } else {
            if len(Sellers) > 1{
                Sellers = Sellers[:j+1]    
            }
            // del Sellers[j+1::]
        }

        // fmt.Println("after delete len(Buyers),len(Sellers) ",len(Buyers),len(Sellers)," ,i,j",i,j)

        pmin := Sellers[j].Prices[t]
        if pmax < pmin {
            pmax = Buyers[i].Prices[t]
        }
        p := pmax
        var buyWinners []string
        var offerPrice []string
        var totalQuantitySeller float64
        var sellWinners []string
        for idx := range Buyers {
            buyWinners = append(buyWinners,Buyers[idx].Identify)
        }
        for idx := range Sellers {
            sellWinners = append(sellWinners,Sellers[idx].Identify)
            offerStr := fmt.Sprintf("%f", Sellers[idx].Prices[t])            
            offerPrice = append(offerPrice,offerStr)
            totalQuantitySeller += Sellers[idx].SurplusEnergy[t]
        }

        pmaxStr := fmt.Sprintf("%f", pmax)
        // pmaxAsBytes,_ := json.Marshal(pmaxStr)

        pminStr := fmt.Sprintf("%f", pmin)
        // pminAsBytes,_ := json.Marshal(pminStr)

        pStr := fmt.Sprintf("%f", p)
        // pAsBytes,_ := json.Marshal(pStr)
        // totalQuantitySellerStr := fmt.Sprintf("%f", totalQuantitySeller)

        var obj trackedModel
        obj.PMax = pmaxStr
        obj.PMin = pminStr
        obj.Price = pStr
        obj.BuyWinners = buyWinners
        obj.SellWinners = sellWinners
        obj.OptimalPrice = "0"
        obj.MaximumSW = 0.0
        obj.OfferPrice = offerPrice
        obj.TotalQuantitySeller = totalQuantitySeller
        obj.TotalQuantityBuyer = 0.0

        objAsBytes,_ := json.Marshal(obj)
        APIstub.PutState(TRACKED_MODEL_KEY, objAsBytes)
        // fmt.Println("error PutState",error)

        var buffer bytes.Buffer
        buffer.WriteString("[")
        buffer.WriteString("{\"winners\":")
        buffer.WriteString("\"['")
        buffer.WriteString(strings.Join(buyWinners, "','"))
        buffer.WriteString("']\"")
        buffer.WriteString(", \"price\":")
        buffer.WriteString(pStr)
        buffer.WriteString("}")
        buffer.WriteString("]")

        fmt.Println("Find the pmax and pmin took ",time.Since(start)," to execute")

        return shim.Success(buffer.Bytes())
    } else {
        var trackedObj trackedModel
        json.Unmarshal(trackedValue, &trackedObj)

        quantityJson := args[1]
        c := make(map[string]string)
        _ = json.Unmarshal([]byte(quantityJson), &c)
        // for key, value := range c {
        //    fmt.Println("key,value",key,value)
        // }
        totalQuantityBuyer := 0.0
        var optimalQuantity []float64
        for _, userID := range trackedObj.BuyWinners {
            // fmt.Println("key",userID,"value",c[userID])
            quantity,_ := strconv.ParseFloat(c[userID], 64)
            totalQuantityBuyer += quantity
            optimalQuantity = append(optimalQuantity,quantity)
        }

        price64,_ := strconv.ParseFloat(trackedObj.Price, 64)
        pmin64,_ := strconv.ParseFloat(trackedObj.PMin, 64)

        totalPriceBenefit := 0.0
        for _,offer := range trackedObj.OfferPrice {
           offer64,_ := strconv.ParseFloat(offer, 64)
           totalPriceBenefit += (price64 - offer64)
        }
        numberOfSeller64 := float64(len(trackedObj.OfferPrice))
        averageSW := totalPriceBenefit/numberOfSeller64 * totalQuantityBuyer

        if averageSW > trackedObj.MaximumSW {
            trackedObj.MaximumSW = averageSW
            trackedObj.OptimalPrice = fmt.Sprintf("%f", price64)
            trackedObj.TotalQuantityBuyer = totalQuantityBuyer
            trackedObj.OptimalQuantity = optimalQuantity
            // fmt.Println("MaximumSW,OptimalPrice,totalQuantityBuyer",trackedObj.MaximumSW,trackedObj.OptimalPrice,totalQuantityBuyer)
        }
        
        // fmt.Println(" The new value of price is ",price64," And pmin is,",pmin64)

        if price64 > pmin64 {
            pStr := fmt.Sprintf("%f", price64 - increment)
            trackedObj.Price = pStr
            objAsBytes,_ := json.Marshal(trackedObj)
            APIstub.PutState(TRACKED_MODEL_KEY, objAsBytes)

            var buffer bytes.Buffer
            buffer.WriteString("[")
            buffer.WriteString("{\"winners\":")
            buffer.WriteString("\"['")
            buffer.WriteString(strings.Join(trackedObj.BuyWinners, "','"))
            buffer.WriteString("']\"")
            buffer.WriteString(", \"price\":")
            buffer.WriteString(pStr)
            buffer.WriteString("}")
            buffer.WriteString("]")

            fmt.Println("Updated price took ",time.Since(start)," to execute. The new value of price is ",price64," And pmin is",pmin64)
            return shim.Success(buffer.Bytes())    
        }else {
            clearPrice,_ := strconv.ParseFloat(trackedObj.OptimalPrice, 64)
            var jsonResponse []ResponseModel
            if trackedObj.TotalQuantitySeller < trackedObj.TotalQuantityBuyer {
                for _, userID := range trackedObj.SellWinners {
                    userAsBytes,_ := APIstub.GetState(userID)
                    user := UserInfo{}
                    json.Unmarshal(userAsBytes, &user)
                    user.ClearPrice[t] = clearPrice
                    user.ClearQuantity[t] = user.SurplusEnergy[t]
                    clearQuantityStr := fmt.Sprintf("%f", user.ClearQuantity[t])
                    res := ResponseModel{UserID:userID,ClearQuantity:clearQuantityStr,ClearPrice:trackedObj.OptimalPrice,Role:"sell"}
                    jsonResponse = append(jsonResponse,res)

                    playerAsBytes, _ := json.Marshal(user)
                    err := APIstub.PutState(userID, playerAsBytes)
                    if err != nil {
                        fmt.Println("err to PutState ",user.Identify," is ", err)        
                    }
                }

                for idx, userID := range trackedObj.BuyWinners {
                    userAsBytes,_ := APIstub.GetState(userID)
                    user := UserInfo{}
                    json.Unmarshal(userAsBytes, &user)
                    user.ClearPrice[t] = clearPrice
                    value :=  (trackedObj.TotalQuantityBuyer - trackedObj.TotalQuantitySeller)*trackedObj.OptimalQuantity[t]/trackedObj.TotalQuantityBuyer
                    user.ClearQuantity[t] = trackedObj.OptimalQuantity[idx] - value
                    clearQuantityStr := fmt.Sprintf("%f", user.ClearQuantity[t])
                    res := ResponseModel{UserID:userID,ClearQuantity:clearQuantityStr,ClearPrice:trackedObj.OptimalPrice,Role:"buy"}
                    jsonResponse = append(jsonResponse,res)

                    playerAsBytes, _ := json.Marshal(user)
                    err := APIstub.PutState(userID, playerAsBytes)
                    if err != nil {
                        fmt.Println("err to PutState ",user.Identify," is ", err)        
                    }
                }
            } else {
                for _, userID := range trackedObj.SellWinners {
                    userAsBytes,_ := APIstub.GetState(userID)
                    user := UserInfo{}
                    json.Unmarshal(userAsBytes, &user)
                    user.ClearPrice[t] = clearPrice
                    value :=  (trackedObj.TotalQuantitySeller - trackedObj.TotalQuantityBuyer)*user.SurplusEnergy[t]/trackedObj.TotalQuantitySeller
                    user.ClearQuantity[t] = user.SurplusEnergy[t] - value
                    clearQuantityStr := fmt.Sprintf("%f", user.ClearQuantity[t])
                    res := ResponseModel{UserID:userID,ClearQuantity:clearQuantityStr,ClearPrice:trackedObj.OptimalPrice,Role:"sell"}
                    jsonResponse = append(jsonResponse,res)

                    playerAsBytes, _ := json.Marshal(user)
                    err := APIstub.PutState(userID, playerAsBytes)
                    if err != nil {
                        fmt.Println("err to PutState ",user.Identify," is ", err)        
                    }
                }

                for idx, userID := range trackedObj.BuyWinners {
                    userAsBytes,_ := APIstub.GetState(userID)
                    user := UserInfo{}
                    json.Unmarshal(userAsBytes, &user)
                    user.ClearPrice[t] = clearPrice
                    user.ClearQuantity[t] = trackedObj.OptimalQuantity[idx]
                    clearQuantityStr := fmt.Sprintf("%f", user.ClearQuantity[t])
                    res := ResponseModel{UserID:userID,ClearQuantity:clearQuantityStr,ClearPrice:trackedObj.OptimalPrice,Role:"buy"}
                    jsonResponse = append(jsonResponse,res)

                    playerAsBytes, _ := json.Marshal(user)
                    err := APIstub.PutState(userID, playerAsBytes)
                    if err != nil {
                        fmt.Println("err to PutState ",user.Identify," is ", err)        
                    }
                }
            }
            // fmt.Println("jsonResponse",jsonResponse)
            result, _ := json.Marshal(jsonResponse)

            return shim.Success(result)             
        }
    }

    return shim.Success(nil)    
    // fmt.Println("user information",Buyers, Sellers) 
}

func (s *SmartContract) emptyContract(APIstub shim.ChaincodeStubInterface) sc.Response {
    contextAsBytes,_ := json.Marshal("context")
    APIstub.PutState("test_PutState", contextAsBytes)
    return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

    // Create a new Smart Contract
    err := shim.Start(new(SmartContract))
    if err != nil {
        fmt.Printf("Error creating new Smart Contract: %s", err)
    }
}
